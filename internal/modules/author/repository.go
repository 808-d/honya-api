package author

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

func (r *Repository) List(ctx context.Context, limit, offset int32) ([]Author, error) {
	rows, err := r.q.ListAuthors(ctx, r.db, sqlc.ListAuthorsParams{PageLimit: limit, PageOffset: offset})
	if err != nil {
		return nil, err
	}
	authors := make([]Author, 0, len(rows))
	for _, row := range rows {
		authors = append(authors, Author{
			ID:         row.ID,
			PublicID:   row.PublicID,
			Name:       row.Name,
			CreateTime: row.CreateTime.Time,
			UpdateTime: row.UpdateTime.Time,
		})
	}
	return authors, nil
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	return r.q.CountAuthors(ctx, r.db)
}

func (r *Repository) Get(ctx context.Context, publicID string) (*Author, error) {
	row, err := r.q.GetAuthorByPublicID(ctx, r.db, publicID)
	if err != nil {
		return nil, err
	}
	g := Author{
		ID:         row.ID,
		PublicID:   row.PublicID,
		Name:       row.Name,
		CreateTime: row.CreateTime.Time,
		UpdateTime: row.UpdateTime.Time,
	}
	return &g, nil
}

func (r *Repository) Create(ctx context.Context, name string) (*Author, error) {
	tx, err := audit.Begin(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	row, err := r.q.CreateAuthor(ctx, tx, name)
	if err != nil {
		return nil, err
	}
	g := Author{
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

func (r *Repository) Update(ctx context.Context, publicID, name string) (*Author, error) {
	tx, err := audit.Begin(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	row, err := r.q.UpdateAuthor(ctx, tx, sqlc.UpdateAuthorParams{Name: name, PublicID: publicID})
	if err != nil {
		return nil, err
	}
	g := Author{
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

	affected, err := r.q.SoftDeleteAuthor(ctx, tx, publicID)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return affected, nil
}
