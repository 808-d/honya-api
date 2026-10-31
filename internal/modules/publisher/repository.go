package publisher

import (
	"context"

	"github.com/808-d/honya-api/internal/database/sqlc"
	"github.com/808-d/honya-api/internal/shared/audit"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db, q: sqlc.New()}
}

func (r *Repository) List(ctx context.Context, limit, offset int32) ([]Publisher, error) {
	rows, err := r.q.ListPublishers(ctx, r.db, sqlc.ListPublishersParams{PageLimit: limit, PageOffset: offset})
	if err != nil {
		return nil, err
	}
	publishers := make([]Publisher, 0, len(rows))
	for _, row := range rows {
		publishers = append(publishers, Publisher{
			ID:         row.ID,
			PublicID:   row.PublicID,
			Name:       row.Name,
			CreateTime: row.CreateTime.Time,
			UpdateTime: row.UpdateTime.Time,
		})
	}
	return publishers, nil
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	return r.q.CountPublishers(ctx, r.db)
}

func (r *Repository) Get(ctx context.Context, publicID string) (*Publisher, error) {
	row, err := r.q.GetPublisherByPublicID(ctx, r.db, publicID)
	if err != nil {
		return nil, err
	}
	g := Publisher{
		ID:         row.ID,
		PublicID:   row.PublicID,
		Name:       row.Name,
		CreateTime: row.CreateTime.Time,
		UpdateTime: row.UpdateTime.Time,
	}
	return &g, nil
}

func (r *Repository) Create(ctx context.Context, name string) (*Publisher, error) {
	tx, err := audit.Begin(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	row, err := r.q.CreatePublisher(ctx, tx, name)
	if err != nil {
		return nil, err
	}
	g := Publisher{
		ID:         row.ID,
		PublicID:   row.PublicID,
		Name:       row.Name,
		CreateTime: row.CreateTime.Time,
		UpdateTime: row.UpdateTime.Time,
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *Repository) Update(ctx context.Context, publicID, name string) (*Publisher, error) {
	tx, err := audit.Begin(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	row, err := r.q.UpdatePublisher(ctx, tx, sqlc.UpdatePublisherParams{Name: name, PublicID: publicID})
	if err != nil {
		return nil, err
	}
	g := Publisher{
		ID:         row.ID,
		PublicID:   row.PublicID,
		Name:       row.Name,
		CreateTime: row.CreateTime.Time,
		UpdateTime: row.UpdateTime.Time,
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *Repository) SoftDelete(ctx context.Context, publicID string) (int64, error) {
	tx, err := audit.Begin(ctx, r.db)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	affected, err := r.q.SoftDeletePublisher(ctx, tx, publicID)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return affected, nil
}
