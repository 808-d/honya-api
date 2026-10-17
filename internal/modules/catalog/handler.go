package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/808-d/honya-api/internal/shared/apperror"
	"github.com/808-d/honya-api/internal/shared/response"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// List handles GET /catalog/products.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filter := ListFilter{
		Status:   optionalInt16(query.Get("status")),
		Type:     optionalInt16(query.Get("type")),
		Search:   query.Get("search"),
		MinPrice: optionalFloat(query.Get("min_price")),
		MaxPrice: optionalFloat(query.Get("max_price")),
		Page:     optionalInt(query.Get("page"), 1),
		PerPage:  optionalInt(query.Get("per_page"), 0),
	}

	result := h.svc.List(r.Context(), filter)
	response.OK(w, "products retrieved", result)
}

// Get handles GET /catalog/products/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	result := h.svc.Get(r.Context(), chi.URLParam(r, "id"))
	response.OK(w, "product retrieved", result)
}

// Create handles POST /catalog/products.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in ProductInput
	if err := decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	result := h.svc.Create(r.Context(), in)
	response.Created(w, "product created", result)
}

// Update handles PUT /catalog/products/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var in ProductInput
	if err := decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	result := h.svc.Update(r.Context(), chi.URLParam(r, "id"), in)
	response.OK(w, "product updated", result)
}

// Delete handles DELETE /catalog/products/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	h.svc.Delete(r.Context(), chi.URLParam(r, "id"))
	response.OK(w, "product deleted", nil)
}

func optionalInt16(raw string) *int16 {
	if raw == "" {
		return nil
	}
	v, err := strconv.ParseInt(raw, 10, 16)
	if err != nil {
		panic(apperror.BadRequest("invalid numeric query parameter"))
	}
	n := int16(v)
	return &n
}

func optionalInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		panic(apperror.BadRequest("invalid numeric query parameter"))
	}
	return v
}

func optionalFloat(raw string) *float64 {
	if raw == "" {
		return nil
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		panic(apperror.BadRequest("invalid numeric query parameter"))
	}
	return &v
}

func decode(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}
