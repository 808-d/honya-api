package product

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/808-d/honya-api/internal/database/sqlc"
	"github.com/808-d/honya-api/internal/shared/apperror"
	"github.com/808-d/honya-api/internal/shared/utils"
	"github.com/jackc/pgx/v5"
)

const (
	defaultPerPage = 20
	maxPerPage     = 100
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List returns a page of products. It panics with an *apperror.Error on failure.
func (s *Service) List(ctx context.Context, filter ListFilter) *ListResult {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > maxPerPage {
		filter.PerPage = defaultPerPage
	}
	search := strings.TrimSpace(filter.Search)
	if filter.MinPrice != nil && filter.MaxPrice != nil && *filter.MinPrice > *filter.MaxPrice {
		panic(apperror.BadRequest("min_price cannot be greater than max_price"))
	}

	products, err := s.repo.List(ctx, sqlc.ListProductsParams{
		FilterStatus: utils.Int2FromPtr(filter.Status),
		FilterType:   utils.Int2FromPtr(filter.Type),
		Search:       utils.TextFromString(search),
		MinPrice:     utils.Float8FromPtr(filter.MinPrice),
		MaxPrice:     utils.Float8FromPtr(filter.MaxPrice),
		PageLimit:    int32(filter.PerPage),
		PageOffset:   int32((filter.Page - 1) * filter.PerPage),
	})
	if err != nil {
		panic(err)
	}

	total, err := s.repo.Count(ctx, sqlc.CountProductsParams{
		FilterStatus: utils.Int2FromPtr(filter.Status),
		FilterType:   utils.Int2FromPtr(filter.Type),
		Search:       utils.TextFromString(search),
		MinPrice:     utils.Float8FromPtr(filter.MinPrice),
		MaxPrice:     utils.Float8FromPtr(filter.MaxPrice),
	})
	if err != nil {
		panic(err)
	}

	items := make([]ProductResponse, 0, len(products))
	for i := range products {
		items = append(items, toResponse(&products[i]))
	}

	return &ListResult{
		Items:   items,
		Total:   total,
		Page:    filter.Page,
		PerPage: filter.PerPage,
	}
}

// Get returns a product by its public id.
func (s *Service) Get(ctx context.Context, publicID string) ProductResponse {
	product, err := s.repo.GetByPublicID(ctx, publicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("product not found"))
		}
		panic(err)
	}
	return toResponse(product)
}

// Create validates and inserts a product.
func (s *Service) Create(ctx context.Context, in ProductInput) ProductResponse {
	product, err := s.repo.Create(ctx, s.validate(in))
	if err != nil {
		panic(err)
	}
	return toResponse(product)
}

// Update validates and replaces a product.
func (s *Service) Update(ctx context.Context, publicID string, in ProductInput) ProductResponse {
	product, err := s.repo.Update(ctx, publicID, s.validate(in))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("product not found"))
		}
		panic(err)
	}
	return toResponse(product)
}

// Delete soft-deletes a product.
func (s *Service) Delete(ctx context.Context, publicID string) {
	affected, err := s.repo.SoftDelete(ctx, publicID)
	if err != nil {
		panic(err)
	}
	if affected == 0 {
		panic(apperror.NotFound("product not found"))
	}
}

func (s *Service) validate(in ProductInput) ProductWrite {
	title := strings.TrimSpace(in.Title)
	if title == "" {
		panic(apperror.BadRequest("title is required"))
	}
	status := ProductStatus(in.Status)
	if !status.Valid() {
		panic(apperror.BadRequest("invalid product status"))
	}
	productType := ProductType(in.Type)
	if !productType.Valid() {
		panic(apperror.BadRequest("invalid product type"))
	}
	if in.Stock < 0 {
		panic(apperror.BadRequest("stock cannot be negative"))
	}
	if in.Price < 0 {
		panic(apperror.BadRequest("price cannot be negative"))
	}

	publishTime, err := utils.ParseTimePtr(in.PublishTime)
	if err != nil {
		panic(apperror.BadRequest("publish_time must be an RFC3339 timestamp"))
	}

	return ProductWrite{
		Title:       title,
		Cover:       in.Cover,
		Description: in.Description,
		PublisherID: in.PublisherID,
		PublishTime: publishTime,
		Stock:       in.Stock,
		Price:       in.Price,
		Status:      status,
		Type:        productType,
		AuthorIDs:   in.AuthorIDs,
		GenreIDs:    in.GenreIDs,
		TagIDs:      in.TagIDs,
	}
}

func toResponse(p *Product) ProductResponse {
	authors := p.AuthorIDs
	if authors == nil {
		authors = []string{}
	}
	genres := p.GenreIDs
	if genres == nil {
		genres = []string{}
	}
	tags := p.TagIDs
	if tags == nil {
		tags = []string{}
	}

	return ProductResponse{
		ID:          p.PublicID,
		Title:       p.Title,
		Cover:       p.Cover,
		Description: p.Description,
		PublisherID: p.PublisherID,
		PublishTime: utils.FormatTimePtr(p.PublishTime),
		Stock:       p.Stock,
		Price:       p.Price,
		Status:      int16(p.Status),
		Type:        int16(p.Type),
		AuthorIDs:   authors,
		GenreIDs:    genres,
		TagIDs:      tags,
		CreateTime:  p.CreateTime.Format(time.RFC3339),
		UpdateTime:  p.UpdateTime.Format(time.RFC3339),
	}
}
