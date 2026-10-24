package coupon

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/808-d/honya-api/internal/database/sqlc"
	"github.com/808-d/honya-api/internal/shared/apperror"
	"github.com/808-d/honya-api/internal/shared/utils"
	"github.com/jackc/pgx/v5"
)

const (
	defaultPerPage = 20
	maxPerPage     = 100

	typeFixed   int16 = 0
	typePercent int16 = 1
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List returns a page of coupons.
func (s *Service) List(ctx context.Context, filter ListFilter) *ListResult {
	page, perPage := normalizePage(filter.Page, filter.PerPage)
	search := strings.TrimSpace(filter.Search)

	coupons, err := s.repo.List(ctx, sqlc.ListCouponsParams{
		Search:       utils.TextFromString(search),
		FilterType:   utils.Int2FromPtr(filter.Type),
		FilterActive: utils.BoolFromPtr(filter.IsActive),
		ActiveNow:    filter.ActiveNow,
		PageOffset:   int32((page - 1) * perPage),
		PageLimit:    int32(perPage),
	})
	if err != nil {
		panic(err)
	}

	total, err := s.repo.Count(ctx, sqlc.CountCouponsParams{
		Search:       utils.TextFromString(search),
		FilterType:   utils.Int2FromPtr(filter.Type),
		FilterActive: utils.BoolFromPtr(filter.IsActive),
		ActiveNow:    filter.ActiveNow,
	})
	if err != nil {
		panic(err)
	}

	items := make([]CouponResponse, 0, len(coupons))
	for i := range coupons {
		items = append(items, toResponse(&coupons[i]))
	}
	return &ListResult{Items: items, Total: total, Page: page, PerPage: perPage}
}

// Get returns a coupon by its public id.
func (s *Service) Get(ctx context.Context, publicID string) CouponResponse {
	c, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("coupon not found"))
		}
		panic(err)
	}
	return toResponse(c)
}

// Create validates and inserts a coupon.
func (s *Service) Create(ctx context.Context, in CouponInput) CouponResponse {
	c, err := s.repo.Create(ctx, s.validateInput(in))
	if err != nil {
		if utils.IsUniqueViolation(err) {
			panic(apperror.Conflict("coupon code already exists"))
		}
		panic(err)
	}
	return toResponse(c)
}

// Update validates and replaces a coupon.
func (s *Service) Update(ctx context.Context, publicID string, in CouponInput) CouponResponse {
	existing, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("coupon not found"))
		}
		panic(err)
	}

	write := s.validateInput(in)
	if write.UsageLimit != nil && *write.UsageLimit < existing.UsedCount {
		panic(apperror.BadRequest("usage_limit cannot be lower than the number of times the coupon has been used"))
	}

	c, err := s.repo.Update(ctx, publicID, write)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("coupon not found"))
		}
		if utils.IsUniqueViolation(err) {
			panic(apperror.Conflict("coupon code already exists"))
		}
		panic(err)
	}
	return toResponse(c)
}

// Delete soft-deletes a coupon.
func (s *Service) Delete(ctx context.Context, publicID string) {
	affected, err := s.repo.SoftDelete(ctx, publicID)
	if err != nil {
		panic(err)
	}
	if affected == 0 {
		panic(apperror.NotFound("coupon not found"))
	}
}

// ListUsages returns a page of redemptions for a coupon.
func (s *Service) ListUsages(ctx context.Context, publicID string, page, perPage int) *UsageListResult {
	c, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("coupon not found"))
		}
		panic(err)
	}

	page, perPage = normalizePage(page, perPage)
	usages, err := s.repo.ListUsages(ctx, c.ID, int32((page-1)*perPage), int32(perPage))
	if err != nil {
		panic(err)
	}
	total, err := s.repo.CountUsages(ctx, c.ID)
	if err != nil {
		panic(err)
	}

	items := make([]UsageResponse, 0, len(usages))
	for _, u := range usages {
		items = append(items, UsageResponse{
			ID:         u.ID,
			UserID:     u.UserID,
			Username:   u.Username,
			OrderID:    u.OrderID,
			CreateTime: u.CreateTime.Format(time.RFC3339),
		})
	}
	return &UsageListResult{Items: items, Total: total, Page: page, PerPage: perPage}
}

// Validate previews a coupon for the current user without consuming it.
func (s *Service) Validate(ctx context.Context, in ValidateRequest, userPublicID string) *ValidateResponse {
	code := normalizeCode(in.Code)
	if code == "" {
		panic(apperror.BadRequest("code is required"))
	}
	if in.Subtotal < 0 {
		panic(apperror.BadRequest("subtotal cannot be negative"))
	}

	userID, err := s.repo.UserIDByPublicID(ctx, userPublicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.Unauthorized("user not found"))
		}
		panic(err)
	}

	c, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("coupon not found"))
		}
		panic(err)
	}

	discount := s.check(ctx, s.repo.db, c, in.Subtotal, userID)
	return &ValidateResponse{
		Code:           c.Code,
		Subtotal:       in.Subtotal,
		DiscountAmount: discount,
		Total:          utils.Round2(in.Subtotal - discount),
	}
}

// Apply validates and locks a coupon inside an order transaction, returning the
// coupon and the discount amount. An empty code yields (nil, 0).
func (s *Service) Apply(ctx context.Context, db sqlc.DBTX, code string, subtotal float64, userID int64) (*Coupon, float64) {
	code = normalizeCode(code)
	if code == "" {
		return nil, 0
	}

	c, err := s.repo.GetByCodeForUpdate(ctx, db, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("coupon not found"))
		}
		panic(err)
	}

	return c, s.check(ctx, db, c, subtotal, userID)
}

// Redeem consumes one use of a coupon for an order.
func (s *Service) Redeem(ctx context.Context, db sqlc.DBTX, couponID, userID, orderID int64) {
	affected, err := s.repo.IncrementUsed(ctx, db, couponID)
	if err != nil {
		panic(err)
	}
	if affected == 0 {
		panic(apperror.Conflict("coupon usage limit reached"))
	}
	if err := s.repo.InsertUsage(ctx, db, couponID, userID, orderID); err != nil {
		panic(err)
	}
}

// Release returns a coupon use when an order is cancelled. It is a no-op when
// the order never used a coupon.
func (s *Service) Release(ctx context.Context, db sqlc.DBTX, orderID int64) {
	couponID, err := s.repo.DeleteUsageByOrder(ctx, db, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return
		}
		panic(err)
	}
	if err := s.repo.DecrementUsed(ctx, db, couponID); err != nil {
		panic(err)
	}
}

// Discount computes the discount for a coupon on a subtotal, rounded to two
// decimals, capped by max_discount and by the subtotal itself.
func Discount(c *Coupon, subtotal float64) float64 {
	if c == nil || subtotal <= 0 {
		return 0
	}

	var discount float64
	switch c.Type {
	case typeFixed:
		discount = c.Value
	case typePercent:
		discount = subtotal * c.Value / 100
		if c.MaxDiscount != nil && discount > *c.MaxDiscount {
			discount = *c.MaxDiscount
		}
	default:
		return 0
	}

	if discount > subtotal {
		discount = subtotal
	}
	if discount < 0 {
		discount = 0
	}
	return utils.Round2(discount)
}

func (s *Service) check(ctx context.Context, db sqlc.DBTX, c *Coupon, subtotal float64, userID int64) float64 {
	if !c.IsActive {
		panic(apperror.BadRequest("coupon is not active"))
	}

	now := time.Now()
	if c.StartAt != nil && now.Before(*c.StartAt) {
		panic(apperror.BadRequest("coupon is not active yet"))
	}
	if c.ExpireAt != nil && now.After(*c.ExpireAt) {
		panic(apperror.BadRequest("coupon has expired"))
	}
	if subtotal < c.MinOrderAmount {
		panic(apperror.BadRequest("order does not meet the minimum amount for this coupon"))
	}
	if c.UsageLimit != nil && c.UsedCount >= *c.UsageLimit {
		panic(apperror.Conflict("coupon usage limit reached"))
	}

	used, err := s.repo.CountUsagesByUser(ctx, db, c.ID, userID)
	if err != nil {
		panic(err)
	}
	if used >= int64(c.PerUserLimit) {
		panic(apperror.Conflict("you have already used this coupon"))
	}

	return Discount(c, subtotal)
}

func (s *Service) validateInput(in CouponInput) CouponWrite {
	code := normalizeCode(in.Code)
	if code == "" {
		panic(apperror.BadRequest("code is required"))
	}
	if in.Type != typeFixed && in.Type != typePercent {
		panic(apperror.BadRequest("type must be 0 (fixed) or 1 (percent)"))
	}
	if in.Value < 0 {
		panic(apperror.BadRequest("value cannot be negative"))
	}
	if in.Type == typePercent && in.Value > 100 {
		panic(apperror.BadRequest("percent value cannot exceed 100"))
	}
	if in.MinOrderAmount < 0 {
		panic(apperror.BadRequest("min_order_amount cannot be negative"))
	}

	var maxDiscount *float64
	if in.Type == typePercent && in.MaxDiscount != nil {
		if *in.MaxDiscount < 0 {
			panic(apperror.BadRequest("max_discount cannot be negative"))
		}
		maxDiscount = in.MaxDiscount
	}

	perUserLimit := int32(1)
	if in.PerUserLimit != nil {
		if *in.PerUserLimit < 1 {
			panic(apperror.BadRequest("per_user_limit must be at least 1"))
		}
		perUserLimit = *in.PerUserLimit
	}
	if in.UsageLimit != nil && *in.UsageLimit < 1 {
		panic(apperror.BadRequest("usage_limit must be at least 1"))
	}

	isActive := true
	if in.IsActive != nil {
		isActive = *in.IsActive
	}

	startAt, err := utils.ParseTimePtr(in.StartAt)
	if err != nil {
		panic(apperror.BadRequest("start_at must be an RFC3339 timestamp"))
	}
	expireAt, err := utils.ParseTimePtr(in.ExpireAt)
	if err != nil {
		panic(apperror.BadRequest("expire_at must be an RFC3339 timestamp"))
	}
	if startAt != nil && expireAt != nil && !expireAt.After(*startAt) {
		panic(apperror.BadRequest("expire_at must be after start_at"))
	}

	return CouponWrite{
		Code:           code,
		Type:           in.Type,
		Value:          in.Value,
		MinOrderAmount: in.MinOrderAmount,
		MaxDiscount:    maxDiscount,
		StartAt:        startAt,
		ExpireAt:       expireAt,
		UsageLimit:     in.UsageLimit,
		PerUserLimit:   perUserLimit,
		IsActive:       isActive,
	}
}

func toResponse(c *Coupon) CouponResponse {
	return CouponResponse{
		ID:             c.PublicID,
		Code:           c.Code,
		Type:           c.Type,
		Value:          c.Value,
		MinOrderAmount: c.MinOrderAmount,
		MaxDiscount:    c.MaxDiscount,
		StartAt:        utils.FormatTimePtr(c.StartAt),
		ExpireAt:       utils.FormatTimePtr(c.ExpireAt),
		UsageLimit:     c.UsageLimit,
		UsedCount:      c.UsedCount,
		PerUserLimit:   c.PerUserLimit,
		IsActive:       c.IsActive,
		CreateTime:     c.CreateTime.Format(time.RFC3339),
		UpdateTime:     c.UpdateTime.Format(time.RFC3339),
	}
}

func normalizeCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
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
