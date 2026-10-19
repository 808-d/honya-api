package cart

import (
	"context"
	"errors"
	"time"

	"github.com/808-d/honya-api/internal/shared/apperror"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List returns the authenticated user's cart.
func (s *Service) List(ctx context.Context, userPublicID string) *CartResponse {
	cartID := s.cartID(ctx, userPublicID)

	items, err := s.repo.ListItems(ctx, cartID)
	if err != nil {
		panic(err)
	}

	result := make([]CartItemResponse, 0, len(items))
	var cartCount int64
	for i := range items {
		result = append(result, toResponse(&items[i]))
		cartCount += int64(items[i].Quantity)
	}
	return &CartResponse{Items: result, CartCount: cartCount}
}

// AddItem adds a product to the cart, incrementing the quantity when the line
// already exists.
func (s *Service) AddItem(ctx context.Context, userPublicID string, in AddItemRequest) *CartMutationResponse {
	if in.ProductID == "" {
		panic(apperror.BadRequest("product_id is required"))
	}
	if in.Quantity < 1 {
		panic(apperror.BadRequest("quantity must be at least 1"))
	}

	productID, err := s.repo.ProductIDByPublicID(ctx, in.ProductID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("product not found"))
		}
		panic(err)
	}

	cartID := s.cartID(ctx, userPublicID)
	item, err := s.repo.AddItem(ctx, cartID, productID, in.Quantity)
	if err != nil {
		panic(err)
	}

	count, err := s.repo.Count(ctx, cartID)
	if err != nil {
		panic(err)
	}
	return &CartMutationResponse{Item: toResponse(item), CartCount: count}
}

// RemoveItem deletes a line from the authenticated user's cart.
func (s *Service) RemoveItem(ctx context.Context, userPublicID string, id int64) *CartMutationResponse {
	cartID := s.cartID(ctx, userPublicID)

	item, err := s.repo.RemoveItem(ctx, id, cartID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("cart item not found"))
		}
		panic(err)
	}

	count, err := s.repo.Count(ctx, cartID)
	if err != nil {
		panic(err)
	}
	return &CartMutationResponse{Item: toResponse(item), CartCount: count}
}

func (s *Service) cartID(ctx context.Context, userPublicID string) int64 {
	cartID, err := s.repo.CartIDByUserPublicID(ctx, userPublicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.Unauthorized("user not found"))
		}
		panic(err)
	}
	return cartID
}

func toResponse(item *CartItem) CartItemResponse {
	return CartItemResponse{
		ID:         item.ID,
		Cover:      item.Cover,
		Title:      item.Title,
		Quantity:   item.Quantity,
		CreateTime: item.CreateTime.Format(time.RFC3339),
	}
}
