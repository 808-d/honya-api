package coupon

import (
	"context"

	"github.com/808-d/honya-api/internal/database/sqlc"
	"github.com/808-d/honya-api/internal/shared/utils"
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

// couponRow mirrors the shared column shape of every coupon query row, so each
// generated row type can be converted to it (struct tags are ignored).
type couponRow struct {
	ID             int64
	PublicID       string
	Code           string
	Type           int16
	Value          float64
	MinOrderAmount float64
	MaxDiscount    pgtype.Numeric
	StartAt        pgtype.Timestamptz
	ExpireAt       pgtype.Timestamptz
	UsageLimit     pgtype.Int4
	UsedCount      int32
	PerUserLimit   int32
	IsActive       bool
	CreateTime     pgtype.Timestamptz
	UpdateTime     pgtype.Timestamptz
}

func toCoupon(r couponRow) *Coupon {
	return &Coupon{
		ID:             r.ID,
		PublicID:       r.PublicID,
		Code:           r.Code,
		Type:           CouponType(r.Type),
		Value:          r.Value,
		MinOrderAmount: r.MinOrderAmount,
		MaxDiscount:    utils.NumericPtr(r.MaxDiscount),
		StartAt:        utils.TimestamptzPtr(r.StartAt),
		ExpireAt:       utils.TimestamptzPtr(r.ExpireAt),
		UsageLimit:     utils.Int4Ptr(r.UsageLimit),
		UsedCount:      r.UsedCount,
		PerUserLimit:   r.PerUserLimit,
		IsActive:       r.IsActive,
		CreateTime:     r.CreateTime.Time,
		UpdateTime:     r.UpdateTime.Time,
	}
}

// Create inserts a new coupon.
func (r *Repository) Create(ctx context.Context, in CouponWrite) (*Coupon, error) {
	row, err := r.q.CreateCoupon(ctx, r.db, sqlc.CreateCouponParams{
		Code:           in.Code,
		Type:           int16(in.Type),
		Value:          in.Value,
		MinOrderAmount: in.MinOrderAmount,
		MaxDiscount:    utils.Float8FromPtr(in.MaxDiscount),
		StartAt:        utils.TimestamptzFromPtr(in.StartAt),
		ExpireAt:       utils.TimestamptzFromPtr(in.ExpireAt),
		UsageLimit:     utils.Int4FromPtr(in.UsageLimit),
		PerUserLimit:   in.PerUserLimit,
		IsActive:       in.IsActive,
	})
	if err != nil {
		return nil, err
	}
	return toCoupon(couponRow(row)), nil
}

// Update replaces an existing coupon.
func (r *Repository) Update(ctx context.Context, publicID string, in CouponWrite) (*Coupon, error) {
	row, err := r.q.UpdateCoupon(ctx, r.db, sqlc.UpdateCouponParams{
		Code:           in.Code,
		Type:           int16(in.Type),
		Value:          in.Value,
		MinOrderAmount: in.MinOrderAmount,
		MaxDiscount:    utils.Float8FromPtr(in.MaxDiscount),
		StartAt:        utils.TimestamptzFromPtr(in.StartAt),
		ExpireAt:       utils.TimestamptzFromPtr(in.ExpireAt),
		UsageLimit:     utils.Int4FromPtr(in.UsageLimit),
		PerUserLimit:   in.PerUserLimit,
		IsActive:       in.IsActive,
		PublicID:       publicID,
	})
	if err != nil {
		return nil, err
	}
	return toCoupon(couponRow(row)), nil
}

// GetByPublicID loads a coupon by its public id.
func (r *Repository) GetByPublicID(ctx context.Context, publicID string) (*Coupon, error) {
	row, err := r.q.GetCouponByPublicID(ctx, r.db, publicID)
	if err != nil {
		return nil, err
	}
	return toCoupon(couponRow(row)), nil
}

// GetByCode loads a coupon by code, case-insensitively.
func (r *Repository) GetByCode(ctx context.Context, code string) (*Coupon, error) {
	row, err := r.q.GetCouponByCode(ctx, r.db, code)
	if err != nil {
		return nil, err
	}
	return toCoupon(couponRow(row)), nil
}

// GetByCodeForUpdate loads and locks a coupon for the rest of the transaction.
func (r *Repository) GetByCodeForUpdate(ctx context.Context, db sqlc.DBTX, code string) (*Coupon, error) {
	row, err := r.q.GetCouponByCodeForUpdate(ctx, db, code)
	if err != nil {
		return nil, err
	}
	return toCoupon(couponRow(row)), nil
}

// List returns a page of coupons.
func (r *Repository) List(ctx context.Context, params sqlc.ListCouponsParams) ([]Coupon, error) {
	rows, err := r.q.ListCoupons(ctx, r.db, params)
	if err != nil {
		return nil, err
	}
	coupons := make([]Coupon, 0, len(rows))
	for _, row := range rows {
		coupons = append(coupons, *toCoupon(couponRow(row)))
	}
	return coupons, nil
}

// Count returns the total number of coupons matching the filters.
func (r *Repository) Count(ctx context.Context, params sqlc.CountCouponsParams) (int64, error) {
	return r.q.CountCoupons(ctx, r.db, params)
}

// SoftDelete marks a coupon as deleted.
func (r *Repository) SoftDelete(ctx context.Context, publicID string) (int64, error) {
	return r.q.SoftDeleteCoupon(ctx, r.db, publicID)
}

// UserIDByPublicID resolves a user's internal id from their public id.
func (r *Repository) UserIDByPublicID(ctx context.Context, publicID string) (int64, error) {
	return r.q.GetUserIDByPublicID(ctx, r.db, publicID)
}

// CountUsagesByUser returns how many times a user has redeemed a coupon.
func (r *Repository) CountUsagesByUser(ctx context.Context, db sqlc.DBTX, couponID, userID int64) (int64, error) {
	return r.q.CountCouponUsagesByUser(ctx, db, sqlc.CountCouponUsagesByUserParams{
		CouponID: couponID,
		UserID:   userID,
	})
}

// IncrementUsed consumes one use, reporting whether it was allowed.
func (r *Repository) IncrementUsed(ctx context.Context, db sqlc.DBTX, couponID int64) (int64, error) {
	return r.q.IncrementCouponUsed(ctx, db, couponID)
}

// DecrementUsed releases one use.
func (r *Repository) DecrementUsed(ctx context.Context, db sqlc.DBTX, couponID int64) error {
	return r.q.DecrementCouponUsed(ctx, db, couponID)
}

// InsertUsage records a redemption.
func (r *Repository) InsertUsage(ctx context.Context, db sqlc.DBTX, couponID, userID, orderID int64) error {
	return r.q.InsertCouponUsage(ctx, db, sqlc.InsertCouponUsageParams{
		CouponID: couponID,
		UserID:   userID,
		OrderID:  orderID,
	})
}

// DeleteUsageByOrder removes the redemption for an order and returns its coupon.
func (r *Repository) DeleteUsageByOrder(ctx context.Context, db sqlc.DBTX, orderID int64) (int64, error) {
	return r.q.DeleteCouponUsageByOrder(ctx, db, orderID)
}

// ListUsages returns a page of redemptions for a coupon.
func (r *Repository) ListUsages(ctx context.Context, couponID int64, offset, limit int32) ([]CouponUsage, error) {
	rows, err := r.q.ListCouponUsages(ctx, r.db, sqlc.ListCouponUsagesParams{
		CouponID:   couponID,
		PageOffset: offset,
		PageLimit:  limit,
	})
	if err != nil {
		return nil, err
	}
	usages := make([]CouponUsage, 0, len(rows))
	for _, row := range rows {
		usages = append(usages, CouponUsage{
			ID:         row.ID,
			UserID:     row.UserID,
			Username:   row.Username,
			OrderID:    row.OrderID,
			CreateTime: row.CreateTime.Time,
		})
	}
	return usages, nil
}

// CountUsages returns the total number of redemptions for a coupon.
func (r *Repository) CountUsages(ctx context.Context, couponID int64) (int64, error) {
	return r.q.CountCouponUsages(ctx, r.db, couponID)
}
