package order

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/808-d/honya-api/internal/database/sqlc"
	"github.com/808-d/honya-api/internal/modules/coupon"
	"github.com/808-d/honya-api/internal/modules/product"
	"github.com/808-d/honya-api/internal/shared/apperror"
	"github.com/808-d/honya-api/internal/shared/utils"
	"github.com/jackc/pgx/v5"
)

const (
	defaultPerPage = 20
	maxPerPage     = 100
)

type Service struct {
	repo               *Repository
	coupons            *coupon.Service
	shippingFeePercent float64
}

func NewService(repo *Repository, coupons *coupon.Service, shippingFeePercent float64) *Service {
	if shippingFeePercent < 0 {
		shippingFeePercent = 0
	}
	return &Service{
		repo:               repo,
		coupons:            coupons,
		shippingFeePercent: shippingFeePercent,
	}
}

// Checkout converts the user's cart into a COD order in a single transaction.
func (s *Service) Checkout(ctx context.Context, userPublicID string, in CreateOrderRequest) *OrderResponse {
	shippingName := strings.TrimSpace(in.ShippingName)
	if shippingName == "" {
		panic(apperror.BadRequest("shipping_name is required"))
	}
	shippingPhone := strings.TrimSpace(in.ShippingPhone)
	if shippingPhone == "" {
		panic(apperror.BadRequest("shipping_phone is required"))
	}
	shippingAddress := strings.TrimSpace(in.ShippingAddress)
	if shippingAddress == "" {
		panic(apperror.BadRequest("shipping_address is required"))
	}

	method := PaymentMethod(in.PaymentMethod)
	if !method.Valid() {
		panic(apperror.BadRequest("only cash on delivery is supported"))
	}

	userID, err := s.repo.UserIDByPublicID(ctx, s.repo.db, userPublicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.Unauthorized("user not found"))
		}
		panic(err)
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	lines, err := s.repo.LockCartForCheckout(ctx, tx, userPublicID)
	if err != nil {
		panic(err)
	}
	if len(lines) == 0 {
		panic(apperror.BadRequest("cart is empty"))
	}

	for _, line := range lines {
		if line.IsDelete {
			panic(apperror.Conflict(fmt.Sprintf("product %q is no longer available", line.ProductTitle)))
		}
		if line.Status != product.ProductStatusActive {
			panic(apperror.Conflict(fmt.Sprintf("product %q is not available", line.ProductTitle)))
		}
		if line.Stock < line.Quantity {
			panic(apperror.Conflict(fmt.Sprintf("insufficient stock for product %q", line.ProductTitle)))
		}
	}

	var subtotal float64
	for _, line := range lines {
		subtotal += utils.Round2(line.UnitPrice * float64(line.Quantity))
	}
	subtotal = utils.Round2(subtotal)

	appliedCoupon, discount := s.coupons.Apply(ctx, tx, in.CouponCode, subtotal, userID)
	shipping, total := s.totals(subtotal, discount)

	var couponID *int64
	if appliedCoupon != nil {
		couponID = &appliedCoupon.ID
	}

	created, err := s.repo.CreateOrder(ctx, tx, CreateOrderWrite{
		UserID:          userID,
		CouponID:        couponID,
		Subtotal:        subtotal,
		DiscountAmount:  discount,
		ShippingFee:     shipping,
		Total:           total,
		Status:          OrderStatusPending,
		PaymentStatus:   PaymentStatusUnpaid,
		PaymentMethod:   method,
		ShippingName:    shippingName,
		ShippingPhone:   shippingPhone,
		ShippingAddress: shippingAddress,
		Note:            strings.TrimSpace(in.Note),
	})
	if err != nil {
		panic(err)
	}

	items := make([]OrderItem, 0, len(lines))
	for _, line := range lines {
		itemSubtotal := utils.Round2(line.UnitPrice * float64(line.Quantity))
		if err := s.repo.CreateOrderItem(ctx, tx, CreateOrderItemWrite{
			OrderID:      created.ID,
			ProductID:    line.ProductID,
			ProductTitle: line.ProductTitle,
			ProductCover: line.ProductCover,
			UnitPrice:    line.UnitPrice,
			Quantity:     line.Quantity,
			Subtotal:     itemSubtotal,
		}); err != nil {
			panic(err)
		}

		productInternalID := line.ProductID
		publicID := line.ProductPublicID
		items = append(items, OrderItem{
			ProductInternalID: &productInternalID,
			ProductID:         &publicID,
			ProductTitle:      line.ProductTitle,
			ProductCover:      line.ProductCover,
			UnitPrice:         line.UnitPrice,
			Quantity:          line.Quantity,
			Subtotal:          itemSubtotal,
		})
	}

	stockLines := append([]CheckoutLine(nil), lines...)
	sort.Slice(stockLines, func(i, j int) bool {
		return stockLines[i].ProductID < stockLines[j].ProductID
	})
	for _, line := range stockLines {
		affected, err := s.repo.DecrementProductStock(ctx, tx, line.ProductID, line.Quantity)
		if err != nil {
			panic(err)
		}
		if affected == 0 {
			panic(apperror.Conflict(fmt.Sprintf("insufficient stock for product %q", line.ProductTitle)))
		}
	}

	if appliedCoupon != nil {
		s.coupons.Redeem(ctx, tx, appliedCoupon.ID, userID, created.ID)
	}

	if err := s.repo.DeleteCartItems(ctx, tx, lines[0].CartID); err != nil {
		panic(err)
	}

	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}

	created.Items = items
	if appliedCoupon != nil {
		created.CouponCode = appliedCoupon.Code
	}
	return toResponse(created, true)
}

// GetForUser returns one of the caller's orders with its items.
func (s *Service) GetForUser(ctx context.Context, userPublicID, publicID string) *OrderResponse {
	o, err := s.repo.GetByPublicIDForUser(ctx, publicID, userPublicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("order not found"))
		}
		panic(err)
	}
	s.loadItems(ctx, s.repo.db, o)
	return toResponse(o, true)
}

// GetAdmin returns any order with its items.
func (s *Service) GetAdmin(ctx context.Context, publicID string) *OrderResponse {
	o, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("order not found"))
		}
		panic(err)
	}
	s.loadItems(ctx, s.repo.db, o)
	return toResponse(o, true)
}

// ListForUser returns a page of the caller's orders without items.
func (s *Service) ListForUser(ctx context.Context, userPublicID string, filter ListFilter) *ListResult {
	page, perPage := normalizePage(filter.Page, filter.PerPage)
	validateStatusFilter(filter.Status)

	orders, err := s.repo.ListForUser(ctx, userPublicID, filter.Status, int32((page-1)*perPage), int32(perPage))
	if err != nil {
		panic(err)
	}
	total, err := s.repo.CountForUser(ctx, userPublicID, filter.Status)
	if err != nil {
		panic(err)
	}
	return listResult(orders, total, page, perPage)
}

// ListAdmin returns a page of all orders without items.
func (s *Service) ListAdmin(ctx context.Context, filter ListFilter) *ListResult {
	page, perPage := normalizePage(filter.Page, filter.PerPage)
	validateStatusFilter(filter.Status)
	validatePaymentStatusFilter(filter.PaymentStatus)

	orders, err := s.repo.ListAdmin(ctx, filter.Status, filter.PaymentStatus, int32((page-1)*perPage), int32(perPage))
	if err != nil {
		panic(err)
	}
	total, err := s.repo.CountAdmin(ctx, filter.Status, filter.PaymentStatus)
	if err != nil {
		panic(err)
	}
	return listResult(orders, total, page, perPage)
}

// Cancel cancels one of the caller's pending, unpaid orders.
func (s *Service) Cancel(ctx context.Context, userPublicID, publicID string) *OrderResponse {
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	o, err := s.repo.GetForUpdateForUser(ctx, tx, publicID, userPublicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("order not found"))
		}
		panic(err)
	}
	if o.Status != OrderStatusPending || o.PaymentStatus != PaymentStatusUnpaid {
		panic(apperror.Conflict("order can only be cancelled while pending and unpaid"))
	}

	s.cancel(ctx, tx, o)

	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}
	return s.GetForUser(ctx, userPublicID, publicID)
}

// UpdateStatus applies an admin fulfillment transition.
func (s *Service) UpdateStatus(ctx context.Context, publicID string, in UpdateStatusRequest) *OrderResponse {
	next := OrderStatus(in.Status)
	if !next.Valid() {
		panic(apperror.BadRequest("invalid order status"))
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	o, err := s.repo.GetForUpdate(ctx, tx, publicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("order not found"))
		}
		panic(err)
	}
	if !o.Status.CanTransitionTo(next) {
		panic(apperror.Conflict("invalid order status transition"))
	}

	switch next {
	case OrderStatusCancelled:
		s.cancel(ctx, tx, o)
	case OrderStatusRefunded:
		s.refund(ctx, tx, o)
	default:
		affected, err := s.repo.UpdateStatus(ctx, tx, o.ID, o.Status, next)
		if err != nil {
			panic(err)
		}
		if affected == 0 {
			panic(apperror.Conflict("order status changed"))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}
	return s.GetAdmin(ctx, publicID)
}

// UpdatePayment applies an admin payment transition.
func (s *Service) UpdatePayment(ctx context.Context, publicID string, in UpdatePaymentRequest) *OrderResponse {
	next := PaymentStatus(in.PaymentStatus)
	if !next.Valid() {
		panic(apperror.BadRequest("invalid payment status"))
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	o, err := s.repo.GetForUpdate(ctx, tx, publicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("order not found"))
		}
		panic(err)
	}
	if !o.PaymentStatus.CanTransitionTo(next) {
		panic(apperror.Conflict("invalid payment status transition"))
	}

	if next == PaymentStatusRefunded {
		s.refund(ctx, tx, o)
	} else {
		var paidAt *time.Time
		if next == PaymentStatusPaid {
			now := time.Now()
			paidAt = &now
		}
		affected, err := s.repo.UpdatePayment(ctx, tx, o.ID, o.PaymentStatus, next, paidAt)
		if err != nil {
			panic(err)
		}
		if affected == 0 {
			panic(apperror.Conflict("payment status changed"))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}
	return s.GetAdmin(ctx, publicID)
}

func (s *Service) cancel(ctx context.Context, tx pgx.Tx, o *Order) {
	affected, err := s.repo.UpdateStatus(ctx, tx, o.ID, o.Status, OrderStatusCancelled)
	if err != nil {
		panic(err)
	}
	if affected == 0 {
		panic(apperror.Conflict("order status changed"))
	}

	s.coupons.Release(ctx, tx, o.ID)

	items, err := s.repo.ListItems(ctx, tx, o.ID)
	if err != nil {
		panic(err)
	}
	sort.SliceStable(items, func(i, j int) bool {
		left, right := items[i].ProductInternalID, items[j].ProductInternalID
		if left == nil {
			return false
		}
		if right == nil {
			return true
		}
		return *left < *right
	})
	for _, item := range items {
		if item.ProductInternalID == nil {
			continue
		}
		if err := s.repo.RestoreProductStock(ctx, tx, *item.ProductInternalID, item.Quantity); err != nil {
			panic(err)
		}
	}
}

func (s *Service) refund(ctx context.Context, tx pgx.Tx, o *Order) {
	if o.Status != OrderStatusCompleted {
		panic(apperror.Conflict("only completed orders can be refunded"))
	}
	if o.PaymentStatus != PaymentStatusPaid {
		panic(apperror.Conflict("order is not paid"))
	}

	affected, err := s.repo.UpdatePayment(ctx, tx, o.ID, o.PaymentStatus, PaymentStatusRefunded, nil)
	if err != nil {
		panic(err)
	}
	if affected == 0 {
		panic(apperror.Conflict("payment status changed"))
	}

	affected, err = s.repo.UpdateStatus(ctx, tx, o.ID, o.Status, OrderStatusRefunded)
	if err != nil {
		panic(err)
	}
	if affected == 0 {
		panic(apperror.Conflict("order status changed"))
	}
}

func (s *Service) loadItems(ctx context.Context, db sqlc.DBTX, o *Order) {
	items, err := s.repo.ListItems(ctx, db, o.ID)
	if err != nil {
		panic(err)
	}
	o.Items = items
}

func (s *Service) totals(subtotal, discount float64) (shipping, total float64) {
	shipping = utils.Round2(subtotal * s.shippingFeePercent / 100)
	total = utils.Round2(subtotal - discount + shipping)
	if total < 0 {
		total = 0
	}
	return shipping, total
}

func listResult(orders []Order, total int64, page, perPage int) *ListResult {
	items := make([]OrderResponse, 0, len(orders))
	for i := range orders {
		items = append(items, *toResponse(&orders[i], false))
	}
	return &ListResult{Items: items, Total: total, Page: page, PerPage: perPage}
}

func toResponse(o *Order, withItems bool) *OrderResponse {
	resp := &OrderResponse{
		ID:              o.PublicID,
		Status:          int16(o.Status),
		PaymentMethod:   int16(o.PaymentMethod),
		PaymentStatus:   int16(o.PaymentStatus),
		Subtotal:        o.Subtotal,
		DiscountAmount:  o.DiscountAmount,
		ShippingFee:     o.ShippingFee,
		Total:           o.Total,
		ShippingName:    o.ShippingName,
		ShippingPhone:   o.ShippingPhone,
		ShippingAddress: o.ShippingAddress,
		Note:            o.Note,
		CouponCode:      o.CouponCode,
		PaidAt:          utils.FormatTimePtr(o.PaidAt),
		CreateTime:      o.CreateTime.Format(time.RFC3339),
		UpdateTime:      o.UpdateTime.Format(time.RFC3339),
	}
	if withItems {
		resp.Items = make([]OrderItemResponse, 0, len(o.Items))
		for _, item := range o.Items {
			resp.Items = append(resp.Items, OrderItemResponse{
				ProductID:    item.ProductID,
				ProductTitle: item.ProductTitle,
				ProductCover: item.ProductCover,
				UnitPrice:    item.UnitPrice,
				Quantity:     item.Quantity,
				Subtotal:     item.Subtotal,
			})
		}
	}
	return resp
}

func validateStatusFilter(status *int16) {
	if status == nil {
		return
	}
	if !OrderStatus(*status).Valid() {
		panic(apperror.BadRequest("invalid status"))
	}
}

func validatePaymentStatusFilter(status *int16) {
	if status == nil {
		return
	}
	if !PaymentStatus(*status).Valid() {
		panic(apperror.BadRequest("invalid payment_status"))
	}
}

func normalizePage(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	return page, perPage
}
