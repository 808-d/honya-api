package dashboard

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts the dashboard stats endpoint behind the access guard.
func (h *Handler) RegisterRoutes(
	r chi.Router,
	auth func(http.Handler) http.Handler,
	access func(http.Handler) http.Handler,
) {
	r.With(auth, access).Get("/admin/stats", h.Stats)
}
