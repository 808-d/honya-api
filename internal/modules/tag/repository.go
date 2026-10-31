package tag

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

func (r *Repository) List(ctx context.Context, limit, offset int32) ([]Tag, error) {
	rows, err := r.q.ListTags(ctx, r.db, sqlc.ListTagsParams{PageLimit: limit, PageOffset: offset})
	if err != nil {
		return nil, err
	}
	tags := make([]Tag, 0, len(rows))
	for _, row := range rows {
		tags = append(tags, Tag{
			ID:         row.ID,
			PublicID:   row.PublicID,
			Name:       row.Name,
			CreateTime: row.CreateTime.Time,
			UpdateTime: row.UpdateTime.Time,
		})
	}
	return tags, nil
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	return r.q.CountTags(ctx, r.db)
}

func (r *Repository) Get(ctx context.Context, publicID string) (*Tag, error) {
	row, err := r.q.GetTagByPublicID(ctx, r.db, publicID)
	if err != nil {
		return nil, err
	}
	g := Tag{
		ID:         row.ID,
		PublicID:   row.PublicID,
		Name:       row.Name,
		CreateTime: row.CreateTime.Time,
		UpdateTime: row.UpdateTime.Time,
	}
	return &g, nil
}

func (r *Repository) Create(ctx context.Context, name string) (*Tag, error) {
	tx, err := audit.Begin(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	row, err := r.q.CreateTag(ctx, tx, name)
	if err != nil {
		return nil, err
	}
	g := Tag{
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

func (r *Repository) Update(ctx context.Context, publicID, name string) (*Tag, error) {
	tx, err := audit.Begin(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	row, err := r.q.UpdateTag(ctx, tx, sqlc.UpdateTagParams{Name: name, PublicID: publicID})
	if err != nil {
		return nil, err
	}
	g := Tag{
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

	affected, err := r.q.SoftDeleteTag(ctx, tx, publicID)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return affected, nil
}
