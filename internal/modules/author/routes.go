package author

import "github.com/go-chi/chi/v5"

// RegisterRoutes mounts the author endpoints on the given router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/authors", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}
