package product

import (
	"net/http"

	"github.com/808-d/honya-api/internal/shared/apperror"
	"github.com/808-d/honya-api/internal/shared/httpx"
	"github.com/808-d/honya-api/internal/shared/response"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// List handles GET /products.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filter := ListFilter{
		Status:   httpx.QueryInt16(query.Get("status")),
		Type:     httpx.QueryInt16(query.Get("type")),
		Search:   query.Get("search"),
		MinPrice: httpx.QueryFloat(query.Get("min_price")),
		MaxPrice: httpx.QueryFloat(query.Get("max_price")),
		Page:     httpx.QueryInt(query.Get("page"), 1),
		PerPage:  httpx.QueryInt(query.Get("per_page"), 0),
	}

	result := h.svc.List(r.Context(), filter)
	response.OK(w, "products retrieved", result)
}

// Get handles GET /products/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	result := h.svc.Get(r.Context(), chi.URLParam(r, "id"))
	response.OK(w, "product retrieved", result)
}

// Create handles POST /products.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in ProductInput
	if err := httpx.Decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	result := h.svc.Create(r.Context(), in)
	response.Created(w, "product created", result)
}

// Update handles PUT /products/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var in ProductInput
	if err := httpx.Decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	result := h.svc.Update(r.Context(), chi.URLParam(r, "id"), in)
	response.OK(w, "product updated", result)
}

// Delete handles DELETE /products/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	h.svc.Delete(r.Context(), chi.URLParam(r, "id"))
	response.OK(w, "product deleted", nil)
}
