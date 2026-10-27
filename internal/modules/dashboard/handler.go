package dashboard

import (
	"net/http"

	"github.com/808-d/honya-api/internal/shared/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Stats handles GET /admin/stats.
func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	response.OK(w, "stats retrieved", h.svc.Stats(r.Context()))
}
