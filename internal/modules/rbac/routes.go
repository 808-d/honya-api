package rbac

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts the admin RBAC endpoints. auth is required everywhere;
// read/write guard the role and user management endpoints by feature.
func (h *Handler) RegisterRoutes(
	r chi.Router,
	auth func(http.Handler) http.Handler,
	roleRead func(http.Handler) http.Handler,
	roleWrite func(http.Handler) http.Handler,
	userRead func(http.Handler) http.Handler,
	userWrite func(http.Handler) http.Handler,
) {
	r.Route("/admin/roles", func(r chi.Router) {
		r.With(auth, roleRead).Get("/", h.ListRoles)
		r.With(auth, roleWrite).Post("/", h.CreateRole)
		r.With(auth, roleRead).Get("/{id}", h.GetRole)
		r.With(auth, roleWrite).Put("/{id}", h.UpdateRole)
		r.With(auth, roleWrite).Delete("/{id}", h.DeleteRole)
	})

	r.With(auth, roleRead).Get("/admin/features", h.ListFeatures)

	r.Route("/admin/users", func(r chi.Router) {
		r.With(auth, userRead).Get("/", h.ListUsers)
		r.With(auth, userWrite).Put("/{id}/roles", h.UpdateUserRoles)
	})
}
