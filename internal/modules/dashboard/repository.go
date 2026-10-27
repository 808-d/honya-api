package dashboard

import (
	"context"

	"github.com/808-d/honya-api/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db, q: sqlc.New()}
}

func (r *Repository) ProductStats(ctx context.Context, lowStock int32) (sqlc.ProductStatsRow, error) {
	return r.q.ProductStats(ctx, r.db, lowStock)
}

func (r *Repository) OrderStats(ctx context.Context) (sqlc.OrderStatsRow, error) {
	return r.q.OrderStats(ctx, r.db)
}

func (r *Repository) OrderStatusCounts(ctx context.Context) ([]sqlc.OrderStatusCountsRow, error) {
	return r.q.OrderStatusCounts(ctx, r.db)
}

func (r *Repository) CouponStats(ctx context.Context) (sqlc.CouponStatsRow, error) {
	return r.q.CouponStats(ctx, r.db)
}

func (r *Repository) UserStats(ctx context.Context) (int64, error) {
	return r.q.UserStats(ctx, r.db)
}
