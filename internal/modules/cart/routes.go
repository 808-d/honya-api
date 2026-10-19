package cart

import "github.com/go-chi/chi/v5"

// RegisterRoutes mounts the cart endpoints on the given router. The caller must
// wrap it with the auth middleware so the user id is available in the context.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/cart", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/items", h.AddItem)
		r.Delete("/items/{id}", h.RemoveItem)
	})
}
