package genre

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

func (r *Repository) List(ctx context.Context, limit, offset int32) ([]Genre, error) {
	rows, err := r.q.ListGenres(ctx, r.db, sqlc.ListGenresParams{PageLimit: limit, PageOffset: offset})
	if err != nil {
		return nil, err
	}
	genres := make([]Genre, 0, len(rows))
	for _, row := range rows {
		genres = append(genres, Genre{
			ID:         row.ID,
			PublicID:   row.PublicID,
			Name:       row.Name,
			CreateTime: row.CreateTime.Time,
			UpdateTime: row.UpdateTime.Time,
		})
	}
	return genres, nil
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	return r.q.CountGenres(ctx, r.db)
}

func (r *Repository) Get(ctx context.Context, publicID string) (*Genre, error) {
	row, err := r.q.GetGenreByPublicID(ctx, r.db, publicID)
	if err != nil {
		return nil, err
	}
	g := Genre{
		ID:         row.ID,
		PublicID:   row.PublicID,
		Name:       row.Name,
		CreateTime: row.CreateTime.Time,
		UpdateTime: row.UpdateTime.Time,
	}
	return &g, nil
}

func (r *Repository) Create(ctx context.Context, name string) (*Genre, error) {
	row, err := r.q.CreateGenre(ctx, r.db, name)
	if err != nil {
		return nil, err
	}
	g := Genre{
		ID:         row.ID,
		PublicID:   row.PublicID,
		Name:       row.Name,
		CreateTime: row.CreateTime.Time,
		UpdateTime: row.UpdateTime.Time,
	}
	return &g, nil
}

func (r *Repository) Update(ctx context.Context, publicID, name string) (*Genre, error) {
	row, err := r.q.UpdateGenre(ctx, r.db, sqlc.UpdateGenreParams{Name: name, PublicID: publicID})
	if err != nil {
		return nil, err
	}
	g := Genre{
		ID:         row.ID,
		PublicID:   row.PublicID,
		Name:       row.Name,
		CreateTime: row.CreateTime.Time,
		UpdateTime: row.UpdateTime.Time,
	}
	return &g, nil
}

func (r *Repository) SoftDelete(ctx context.Context, publicID string) (int64, error) {
	return r.q.SoftDeleteGenre(ctx, r.db, publicID)
}
