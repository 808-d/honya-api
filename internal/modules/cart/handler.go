package cart

import (
	"net/http"
	"strconv"

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

// List handles GET /cart.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	result := h.svc.List(r.Context(), middleware.UserIDFromContext(r.Context()))
	response.OK(w, "cart retrieved", result)
}

// AddItem handles POST /cart/items.
func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	var in AddItemRequest
	if err := httpx.Decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	result := h.svc.AddItem(r.Context(), middleware.UserIDFromContext(r.Context()), in)
	response.Created(w, "item added to cart", result)
}

// RemoveItem handles DELETE /cart/items/{id}.
func (h *Handler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		panic(apperror.BadRequest("invalid cart item id"))
	}

	result := h.svc.RemoveItem(r.Context(), middleware.UserIDFromContext(r.Context()), id)
	response.OK(w, "item removed from cart", result)
}
