package coupon

import (
	"net/http"

	"github.com/808-d/honya-api/internal/middleware"
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

// List handles GET /coupons.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	activeNow := httpx.QueryBool(query.Get("active_now"))

	filter := ListFilter{
		Search:    query.Get("search"),
		Type:      httpx.QueryInt16(query.Get("type")),
		IsActive:  httpx.QueryBool(query.Get("is_active")),
		ActiveNow: activeNow != nil && *activeNow,
		Page:      httpx.QueryInt(query.Get("page"), 1),
		PerPage:   httpx.QueryInt(query.Get("per_page"), 0),
	}

	result := h.svc.List(r.Context(), filter)
	response.OK(w, "coupons retrieved", result)
}

// Get handles GET /coupons/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	result := h.svc.Get(r.Context(), chi.URLParam(r, "id"))
	response.OK(w, "coupon retrieved", result)
}

// Create handles POST /coupons.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in CouponInput
	if err := httpx.Decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	result := h.svc.Create(r.Context(), in)
	response.Created(w, "coupon created", result)
}

// Update handles PUT /coupons/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var in CouponInput
	if err := httpx.Decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	result := h.svc.Update(r.Context(), chi.URLParam(r, "id"), in)
	response.OK(w, "coupon updated", result)
}

// Delete handles DELETE /coupons/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	h.svc.Delete(r.Context(), chi.URLParam(r, "id"))
	response.OK(w, "coupon deleted", nil)
}

// ListUsages handles GET /coupons/{id}/usages.
func (h *Handler) ListUsages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	result := h.svc.ListUsages(
		r.Context(),
		chi.URLParam(r, "id"),
		httpx.QueryInt(query.Get("page"), 1),
		httpx.QueryInt(query.Get("per_page"), 0),
	)
	response.OK(w, "coupon usages retrieved", result)
}

// Validate handles POST /coupons/validate for any authenticated user.
func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	var in ValidateRequest
	if err := httpx.Decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	result := h.svc.Validate(r.Context(), in, middleware.UserIDFromContext(r.Context()))
	response.OK(w, "coupon applied", result)
}
