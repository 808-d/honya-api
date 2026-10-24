package coupon

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts the coupon endpoints. auth is required on every route;
// read and write guard the admin endpoints by feature.
func (h *Handler) RegisterRoutes(
	r chi.Router,
	auth func(http.Handler) http.Handler,
	read func(http.Handler) http.Handler,
	write func(http.Handler) http.Handler,
) {
	r.Route("/coupons", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(auth)
			r.Post("/validate", h.Validate)
		})

		r.Group(func(r chi.Router) {
			r.Use(auth, read)
			r.Get("/", h.List)
			r.Get("/{id}", h.Get)
			r.Get("/{id}/usages", h.ListUsages)
		})

		r.Group(func(r chi.Router) {
			r.Use(auth, write)
			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
		})
	})
}
