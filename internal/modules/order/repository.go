package order

import (
	"context"
	"time"

	"github.com/808-d/honya-api/internal/database/sqlc"
	"github.com/808-d/honya-api/internal/modules/product"
	"github.com/808-d/honya-api/internal/shared/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db, q: sqlc.New()}
}

// Begin starts a transaction owned by the service.
func (r *Repository) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}

// orderRow mirrors the shared column shape of every order read query, so each
// generated row type can be converted to it (struct tags are ignored).
type orderRow struct {
	ID              int64
	PublicID        string
	UserID          int64
	CouponCode      string
	Subtotal        float64
	DiscountAmount  float64
	ShippingFee     float64
	Total           float64
	Status          int16
	PaymentStatus   int16
	PaymentMethod   int16
	ShippingName    string
	ShippingPhone   string
	ShippingAddress string
	Note            string
	PaidAt          pgtype.Timestamptz
	CreateTime      pgtype.Timestamptz
	UpdateTime      pgtype.Timestamptz
}

func toOrder(r orderRow) *Order {
	return &Order{
		ID:              r.ID,
		PublicID:        r.PublicID,
		UserID:          r.UserID,
		CouponCode:      r.CouponCode,
		Subtotal:        r.Subtotal,
		DiscountAmount:  r.DiscountAmount,
		ShippingFee:     r.ShippingFee,
		Total:           r.Total,
		Status:          OrderStatus(r.Status),
		PaymentStatus:   PaymentStatus(r.PaymentStatus),
		PaymentMethod:   PaymentMethod(r.PaymentMethod),
		ShippingName:    r.ShippingName,
		ShippingPhone:   r.ShippingPhone,
		ShippingAddress: r.ShippingAddress,
		Note:            r.Note,
		PaidAt:          utils.TimestamptzPtr(r.PaidAt),
		CreateTime:      r.CreateTime.Time,
		UpdateTime:      r.UpdateTime.Time,
	}
}

// UserIDByPublicID resolves a user's internal id from their public id.
func (r *Repository) UserIDByPublicID(ctx context.Context, db sqlc.DBTX, publicID string) (int64, error) {
	return r.q.GetUserIDByPublicID(ctx, db, publicID)
}

// LockCartForCheckout locks the user's cart lines and returns their product
// snapshots for checkout.
func (r *Repository) LockCartForCheckout(ctx context.Context, db sqlc.DBTX, userPublicID string) ([]CheckoutLine, error) {
	rows, err := r.q.LockCartForCheckout(ctx, db, userPublicID)
	if err != nil {
		return nil, err
	}

	lines := make([]CheckoutLine, 0, len(rows))
	for _, row := range rows {
		lines = append(lines, CheckoutLine{
			CartItemID:      row.CartItemID,
			CartID:          row.CartID,
			ProductID:       row.ProductID,
			ProductPublicID: row.ProductPublicID,
			ProductTitle:    row.ProductTitle,
			ProductCover:    row.ProductCover,
			UnitPrice:       row.UnitPrice,
			Quantity:        row.Quantity,
			Stock:           row.Stock,
			Status:          product.ProductStatus(row.Status),
			IsDelete:        row.IsDelete,
		})
	}
	return lines, nil
}

// CreateOrder inserts the order header.
func (r *Repository) CreateOrder(ctx context.Context, db sqlc.DBTX, in CreateOrderWrite) (*Order, error) {
	row, err := r.q.CreateOrder(ctx, db, sqlc.CreateOrderParams{
		UserID:          in.UserID,
		CouponID:        utils.Int8FromPtr(in.CouponID),
		Subtotal:        in.Subtotal,
		DiscountAmount:  in.DiscountAmount,
		ShippingFee:     in.ShippingFee,
		Total:           in.Total,
		Status:          int16(in.Status),
		PaymentStatus:   int16(in.PaymentStatus),
		PaymentMethod:   int16(in.PaymentMethod),
		ShippingName:    utils.TextFromString(in.ShippingName),
		ShippingPhone:   utils.TextFromString(in.ShippingPhone),
		ShippingAddress: utils.TextFromString(in.ShippingAddress),
		Note:            utils.TextFromString(in.Note),
	})
	if err != nil {
		return nil, err
	}

	return &Order{
		ID:              row.ID,
		PublicID:        row.PublicID,
		UserID:          in.UserID,
		Subtotal:        in.Subtotal,
		DiscountAmount:  in.DiscountAmount,
		ShippingFee:     in.ShippingFee,
		Total:           in.Total,
		Status:          in.Status,
		PaymentStatus:   in.PaymentStatus,
		PaymentMethod:   in.PaymentMethod,
		ShippingName:    in.ShippingName,
		ShippingPhone:   in.ShippingPhone,
		ShippingAddress: in.ShippingAddress,
		Note:            in.Note,
		CreateTime:      row.CreateTime.Time,
		UpdateTime:      row.UpdateTime.Time,
	}, nil
}

// CreateOrderItem inserts one immutable order line.
func (r *Repository) CreateOrderItem(ctx context.Context, db sqlc.DBTX, in CreateOrderItemWrite) error {
	productID := in.ProductID
	return r.q.CreateOrderItem(ctx, db, sqlc.CreateOrderItemParams{
		OrderID:      in.OrderID,
		ProductID:    utils.Int8FromPtr(&productID),
		ProductTitle: in.ProductTitle,
		ProductCover: in.ProductCover,
		UnitPrice:    in.UnitPrice,
		Quantity:     in.Quantity,
		Subtotal:     in.Subtotal,
	})
}

// DecrementProductStock subtracts stock only when enough is available.
func (r *Repository) DecrementProductStock(ctx context.Context, db sqlc.DBTX, productID int64, quantity int32) (int64, error) {
	return r.q.DecrementProductStock(ctx, db, sqlc.DecrementProductStockParams{
		ProductID: productID,
		Quantity:  quantity,
	})
}

// RestoreProductStock adds stock back when an order is cancelled.
func (r *Repository) RestoreProductStock(ctx context.Context, db sqlc.DBTX, productID int64, quantity int32) error {
	return r.q.RestoreProductStock(ctx, db, sqlc.RestoreProductStockParams{
		ProductID: productID,
		Quantity:  quantity,
	})
}

// DeleteCartItems clears the cart after a successful checkout.
func (r *Repository) DeleteCartItems(ctx context.Context, db sqlc.DBTX, cartID int64) error {
	return r.q.DeleteCartItems(ctx, db, cartID)
}

// GetByPublicID loads any order by public id.
func (r *Repository) GetByPublicID(ctx context.Context, publicID string) (*Order, error) {
	row, err := r.q.GetOrderByPublicID(ctx, r.db, publicID)
	if err != nil {
		return nil, err
	}
	return toOrder(orderRow(row)), nil
}

// GetByPublicIDForUser loads an order owned by the given user.
func (r *Repository) GetByPublicIDForUser(ctx context.Context, publicID, userPublicID string) (*Order, error) {
	row, err := r.q.GetOrderByPublicIDForUser(ctx, r.db, sqlc.GetOrderByPublicIDForUserParams{
		PublicID:     publicID,
		UserPublicID: userPublicID,
	})
	if err != nil {
		return nil, err
	}
	return toOrder(orderRow(row)), nil
}

// GetForUpdate locks an order for the rest of the transaction.
func (r *Repository) GetForUpdate(ctx context.Context, db sqlc.DBTX, publicID string) (*Order, error) {
	row, err := r.q.GetOrderForUpdate(ctx, db, publicID)
	if err != nil {
		return nil, err
	}
	return toOrder(orderRow(row)), nil
}

// GetForUpdateForUser locks an order owned by the given user.
func (r *Repository) GetForUpdateForUser(ctx context.Context, db sqlc.DBTX, publicID, userPublicID string) (*Order, error) {
	row, err := r.q.GetOrderForUpdateForUser(ctx, db, sqlc.GetOrderForUpdateForUserParams{
		PublicID:     publicID,
		UserPublicID: userPublicID,
	})
	if err != nil {
		return nil, err
	}
	return toOrder(orderRow(row)), nil
}

// ListItems returns the snapshot lines of an order.
func (r *Repository) ListItems(ctx context.Context, db sqlc.DBTX, orderID int64) ([]OrderItem, error) {
	rows, err := r.q.ListOrderItems(ctx, db, orderID)
	if err != nil {
		return nil, err
	}
	items := make([]OrderItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, OrderItem{
			ID:                row.ID,
			ProductInternalID: utils.Int8Ptr(row.ProductInternalID),
			ProductID:         utils.OptionalString(row.ProductID),
			ProductTitle:      row.ProductTitle,
			ProductCover:      row.ProductCover,
			UnitPrice:         row.UnitPrice,
			Quantity:          row.Quantity,
			Subtotal:          row.Subtotal,
		})
	}
	return items, nil
}

// ListForUser returns a page of the user's orders.
func (r *Repository) ListForUser(ctx context.Context, userPublicID string, status *int16, offset, limit int32) ([]Order, error) {
	rows, err := r.q.ListOrdersForUser(ctx, r.db, sqlc.ListOrdersForUserParams{
		UserPublicID: userPublicID,
		FilterStatus: utils.Int2FromPtr(status),
		PageOffset:   offset,
		PageLimit:    limit,
	})
	if err != nil {
		return nil, err
	}
	orders := make([]Order, 0, len(rows))
	for _, row := range rows {
		orders = append(orders, *toOrder(orderRow(row)))
	}
	return orders, nil
}

// CountForUser returns the total number of the user's orders.
func (r *Repository) CountForUser(ctx context.Context, userPublicID string, status *int16) (int64, error) {
	return r.q.CountOrdersForUser(ctx, r.db, sqlc.CountOrdersForUserParams{
		UserPublicID: userPublicID,
		FilterStatus: utils.Int2FromPtr(status),
	})
}

// ListAdmin returns a page of all orders.
func (r *Repository) ListAdmin(ctx context.Context, status, paymentStatus *int16, offset, limit int32) ([]Order, error) {
	rows, err := r.q.ListOrders(ctx, r.db, sqlc.ListOrdersParams{
		FilterStatus:        utils.Int2FromPtr(status),
		FilterPaymentStatus: utils.Int2FromPtr(paymentStatus),
		PageOffset:          offset,
		PageLimit:           limit,
	})
	if err != nil {
		return nil, err
	}
	orders := make([]Order, 0, len(rows))
	for _, row := range rows {
		orders = append(orders, *toOrder(orderRow(row)))
	}
	return orders, nil
}

// CountAdmin returns the total number of orders.
func (r *Repository) CountAdmin(ctx context.Context, status, paymentStatus *int16) (int64, error) {
	return r.q.CountOrders(ctx, r.db, sqlc.CountOrdersParams{
		FilterStatus:        utils.Int2FromPtr(status),
		FilterPaymentStatus: utils.Int2FromPtr(paymentStatus),
	})
}

// UpdateStatus applies a guarded status transition.
func (r *Repository) UpdateStatus(ctx context.Context, db sqlc.DBTX, id int64, expected, next OrderStatus) (int64, error) {
	return r.q.UpdateOrderStatus(ctx, db, sqlc.UpdateOrderStatusParams{
		ID:             id,
		ExpectedStatus: int16(expected),
		Status:         int16(next),
	})
}

// UpdatePayment applies a guarded payment transition.
func (r *Repository) UpdatePayment(ctx context.Context, db sqlc.DBTX, id int64, expected, next PaymentStatus, paidAt *time.Time) (int64, error) {
	return r.q.UpdateOrderPayment(ctx, db, sqlc.UpdateOrderPaymentParams{
		ID:             id,
		ExpectedStatus: int16(expected),
		PaymentStatus:  int16(next),
		PaidAt:         utils.TimestamptzFromPtr(paidAt),
	})
}
