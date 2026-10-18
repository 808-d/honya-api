package author

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/808-d/honya-api/internal/shared/apperror"
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

func (s *Service) List(ctx context.Context, page, perPage int) *ListResult {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}

	authors, err := s.repo.List(ctx, int32(perPage), int32((page-1)*perPage))
	if err != nil {
		panic(err)
	}
	total, err := s.repo.Count(ctx)
	if err != nil {
		panic(err)
	}

	items := make([]Response, 0, len(authors))
	for i := range authors {
		items = append(items, toResponse(&authors[i]))
	}
	return &ListResult{Items: items, Total: total, Page: page, PerPage: perPage}
}

func (s *Service) Get(ctx context.Context, publicID string) Response {
	g, err := s.repo.Get(ctx, publicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("author not found"))
		}
		panic(err)
	}
	return toResponse(g)
}

func (s *Service) Create(ctx context.Context, in Input) Response {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		panic(apperror.BadRequest("name is required"))
	}
	g, err := s.repo.Create(ctx, name)
	if err != nil {
		panic(err)
	}
	return toResponse(g)
}

func (s *Service) Update(ctx context.Context, publicID string, in Input) Response {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		panic(apperror.BadRequest("name is required"))
	}
	g, err := s.repo.Update(ctx, publicID, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.NotFound("author not found"))
		}
		panic(err)
	}
	return toResponse(g)
}

func (s *Service) Delete(ctx context.Context, publicID string) {
	affected, err := s.repo.SoftDelete(ctx, publicID)
	if err != nil {
		panic(err)
	}
	if affected == 0 {
		panic(apperror.NotFound("author not found"))
	}
}

func toResponse(g *Author) Response {
	return Response{
		ID:         g.PublicID,
		Name:       g.Name,
		CreateTime: g.CreateTime.Format(time.RFC3339),
		UpdateTime: g.UpdateTime.Format(time.RFC3339),
	}
}
