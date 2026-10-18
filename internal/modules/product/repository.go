package product

import (
	"context"
	"time"

	"github.com/808-d/honya-api/internal/database/sqlc"
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

func (r *Repository) List(ctx context.Context, params sqlc.ListProductsParams) ([]Product, error) {
	rows, err := r.q.ListProducts(ctx, r.db, params)
	if err != nil {
		return nil, err
	}
	products := make([]Product, 0, len(rows))
	for _, row := range rows {
		products = append(products, Product{
			ID:          row.ID,
			PublicID:    row.PublicID,
			Title:       row.Title,
			Cover:       row.Cover,
			Description: row.Description,
			PublisherID: optionalString(row.PublisherID),
			PublishTime: timePtr(row.PublishTime),
			Stock:       row.Stock,
			Price:       row.Price,
			Status:      row.Status,
			Type:        row.Type,
			CreateTime:  row.CreateTime.Time,
			UpdateTime:  row.UpdateTime.Time,
		})
	}
	return products, nil
}

func (r *Repository) Count(ctx context.Context, params sqlc.CountProductsParams) (int64, error) {
	return r.q.CountProducts(ctx, r.db, params)
}

func (r *Repository) GetByPublicID(ctx context.Context, publicID string) (*Product, error) {
	row, err := r.q.GetProductByPublicID(ctx, r.db, publicID)
	if err != nil {
		return nil, err
	}
	product := &Product{
		ID:          row.ID,
		PublicID:    row.PublicID,
		Title:       row.Title,
		Cover:       row.Cover,
		Description: row.Description,
		PublisherID: optionalString(row.PublisherID),
		PublishTime: timePtr(row.PublishTime),
		Stock:       row.Stock,
		Price:       row.Price,
		Status:      row.Status,
		Type:        row.Type,
		CreateTime:  row.CreateTime.Time,
		UpdateTime:  row.UpdateTime.Time,
	}
	// TODO: using join
	if err := r.loadRelations(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (r *Repository) Create(ctx context.Context, in ProductWrite) (*Product, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	row, err := r.q.CreateProduct(ctx, tx, sqlc.CreateProductParams{
		Title:       in.Title,
		Cover:       in.Cover,
		Description: in.Description,
		PublisherID: textFromPtr(in.PublisherID),
		PublishTime: timestamptzFromPtr(in.PublishTime),
		Stock:       in.Stock,
		Price:       in.Price,
		Status:      in.Status,
		Type:        in.Type,
	})
	if err != nil {
		return nil, err
	}
	if err := r.setRelations(ctx, tx, row.ID, in.AuthorIDs, in.GenreIDs, in.TagIDs); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetByPublicID(ctx, row.PublicID)
}

func (r *Repository) Update(ctx context.Context, publicID string, in ProductWrite) (*Product, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	row, err := r.q.UpdateProduct(ctx, tx, sqlc.UpdateProductParams{
		Title:       in.Title,
		Cover:       in.Cover,
		Description: in.Description,
		PublisherID: textFromPtr(in.PublisherID),
		PublishTime: timestamptzFromPtr(in.PublishTime),
		Stock:       in.Stock,
		Price:       in.Price,
		Status:      in.Status,
		Type:        in.Type,
		PublicID:    publicID,
	})
	if err != nil {
		return nil, err
	}
	if err := r.setRelations(ctx, tx, row.ID, in.AuthorIDs, in.GenreIDs, in.TagIDs); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetByPublicID(ctx, row.PublicID)
}

func (r *Repository) SoftDelete(ctx context.Context, publicID string) (int64, error) {
	return r.q.SoftDeleteProduct(ctx, r.db, publicID)
}

func (r *Repository) loadRelations(ctx context.Context, product *Product) error {
	authors, err := r.q.GetProductAuthorIDs(ctx, r.db, product.ID)
	if err != nil {
		return err
	}
	genres, err := r.q.GetProductGenreIDs(ctx, r.db, product.ID)
	if err != nil {
		return err
	}
	tags, err := r.q.GetProductTagIDs(ctx, r.db, product.ID)
	if err != nil {
		return err
	}
	product.AuthorIDs = authors
	product.GenreIDs = genres
	product.TagIDs = tags
	return nil
}

func (r *Repository) setRelations(ctx context.Context, db sqlc.DBTX, productID int64, authorIDs, genreIDs, tagIDs []string) error {
	if err := r.q.ClearProductAuthors(ctx, db, productID); err != nil {
		return err
	}
	if err := r.q.ClearProductGenres(ctx, db, productID); err != nil {
		return err
	}
	if err := r.q.ClearProductTags(ctx, db, productID); err != nil {
		return err
	}

	if len(authorIDs) > 0 {
		if err := r.q.AddProductAuthors(ctx, db, sqlc.AddProductAuthorsParams{ProductID: productID, PublicIds: authorIDs}); err != nil {
			return err
		}
	}
	if len(genreIDs) > 0 {
		if err := r.q.AddProductGenres(ctx, db, sqlc.AddProductGenresParams{ProductID: productID, PublicIds: genreIDs}); err != nil {
			return err
		}
	}
	if len(tagIDs) > 0 {
		if err := r.q.AddProductTags(ctx, db, sqlc.AddProductTagsParams{ProductID: productID, PublicIds: tagIDs}); err != nil {
			return err
		}
	}
	return nil
}

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func textFromPtr(s *string) pgtype.Text {
	if s == nil || *s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func timestamptzFromPtr(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func timePtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	t := ts.Time
	return &t
}
