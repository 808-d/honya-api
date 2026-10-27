package rbac

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

// ListRoles handles GET /admin/roles.
func (h *Handler) ListRoles(w http.ResponseWriter, r *http.Request) {
	response.OK(w, "roles retrieved", h.svc.ListRoles(r.Context()))
}

// GetRole handles GET /admin/roles/{id}.
func (h *Handler) GetRole(w http.ResponseWriter, r *http.Request) {
	response.OK(w, "role retrieved", h.svc.GetRole(r.Context(), chi.URLParam(r, "id")))
}

// CreateRole handles POST /admin/roles.
func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var in RoleInput
	if err := httpx.Decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}
	response.Created(w, "role created", h.svc.CreateRole(r.Context(), in))
}

// UpdateRole handles PUT /admin/roles/{id}.
func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	var in RoleInput
	if err := httpx.Decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}
	actor := middleware.UserIDFromContext(r.Context())
	response.OK(w, "role updated", h.svc.UpdateRole(r.Context(), actor, chi.URLParam(r, "id"), in))
}

// DeleteRole handles DELETE /admin/roles/{id}.
func (h *Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	actor := middleware.UserIDFromContext(r.Context())
	h.svc.DeleteRole(r.Context(), actor, chi.URLParam(r, "id"))
	response.OK(w, "role deleted", nil)
}

// ListFeatures handles GET /admin/features.
func (h *Handler) ListFeatures(w http.ResponseWriter, r *http.Request) {
	response.OK(w, "features retrieved", h.svc.ListFeatures(r.Context()))
}

// ListUsers handles GET /admin/users.
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	result := h.svc.ListUsers(
		r.Context(),
		query.Get("search"),
		httpx.QueryInt(query.Get("page"), 1),
		httpx.QueryInt(query.Get("per_page"), 0),
	)
	response.OK(w, "users retrieved", result)
}

// UpdateUserRoles handles PUT /admin/users/{id}/roles.
func (h *Handler) UpdateUserRoles(w http.ResponseWriter, r *http.Request) {
	var in UpdateUserRolesRequest
	if err := httpx.Decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}
	actor := middleware.UserIDFromContext(r.Context())
	response.OK(w, "user roles updated", h.svc.UpdateUserRoles(r.Context(), actor, chi.URLParam(r, "id"), in.RoleIDs))
}
