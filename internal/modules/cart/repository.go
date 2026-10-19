package cart

import (
	"context"

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

// CartIDByUserPublicID returns the user's cart id, creating the cart on first use.
func (r *Repository) CartIDByUserPublicID(ctx context.Context, userPublicID string) (int64, error) {
	return r.q.GetOrCreateCart(ctx, r.db, userPublicID)
}

// ProductIDByPublicID resolves a product's internal id from its public id.
func (r *Repository) ProductIDByPublicID(ctx context.Context, publicID string) (int64, error) {
	row, err := r.q.GetProductByPublicID(ctx, r.db, publicID)
	if err != nil {
		return 0, err
	}
	return row.ID, nil
}

// ListItems returns every line in the cart.
func (r *Repository) ListItems(ctx context.Context, cartID int64) ([]CartItem, error) {
	rows, err := r.q.GetCartItems(ctx, r.db, cartID)
	if err != nil {
		return nil, err
	}
	items := make([]CartItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, toItem(row.ID, row.Cover, row.Title, row.Quantity, row.CreateTime))
	}
	return items, nil
}

// Count returns the total quantity of visible lines in the cart.
func (r *Repository) Count(ctx context.Context, cartID int64) (int64, error) {
	return r.q.CountCartItems(ctx, r.db, cartID)
}

// AddItem inserts a line or increments its quantity when it already exists.
func (r *Repository) AddItem(ctx context.Context, cartID, productID int64, quantity int32) (*CartItem, error) {
	row, err := r.q.AddToCart(ctx, r.db, sqlc.AddToCartParams{
		CartID:    cartID,
		ProductID: productID,
		Quantity:  quantity,
	})
	if err != nil {
		return nil, err
	}
	item := toItem(row.ID, row.Cover, row.Title, row.Quantity, row.CreateTime)
	return &item, nil
}

// RemoveItem deletes a line from the cart and returns the removed line.
func (r *Repository) RemoveItem(ctx context.Context, id, cartID int64) (*CartItem, error) {
	row, err := r.q.RemoveFromCart(ctx, r.db, sqlc.RemoveFromCartParams{ID: id, CartID: cartID})
	if err != nil {
		return nil, err
	}
	item := toItem(row.ID, row.Cover, row.Title, row.Quantity, row.CreateTime)
	return &item, nil
}

func toItem(id int64, cover, title string, quantity int32, createTime pgtype.Timestamptz) CartItem {
	return CartItem{
		ID:         id,
		Cover:      cover,
		Title:      title,
		Quantity:   quantity,
		CreateTime: createTime.Time,
	}
}
