package order

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts the customer and admin order endpoints. auth is
// required on every route; read and write guard the admin endpoints by feature.
func (h *Handler) RegisterRoutes(
	r chi.Router,
	auth func(http.Handler) http.Handler,
	read func(http.Handler) http.Handler,
	write func(http.Handler) http.Handler,
) {
	r.Route("/orders", func(r chi.Router) {
		r.Use(auth)
		r.Post("/", h.Checkout)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Post("/{id}/cancel", h.Cancel)
	})

	r.Route("/admin/orders", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(auth, read)
			r.Get("/", h.AdminList)
			r.Get("/{id}", h.AdminGet)
		})

		r.Group(func(r chi.Router) {
			r.Use(auth, write)
			r.Patch("/{id}/status", h.AdminUpdateStatus)
			r.Patch("/{id}/payment", h.AdminUpdatePayment)
		})
	})
}
