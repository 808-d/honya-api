package author

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts the author endpoints on the given router.
func (h *Handler) RegisterRoutes(r chi.Router, auth, write func(http.Handler) http.Handler) {
	r.Route("/authors", func(r chi.Router) {
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.With(auth, write).Post("/", h.Create)
		r.With(auth, write).Put("/{id}", h.Update)
		r.With(auth, write).Delete("/{id}", h.Delete)
	})
}
