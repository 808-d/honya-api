package tag

import (
	"net/http"

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

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	result := h.svc.List(r.Context(), httpx.QueryInt(q.Get("page"), 1), httpx.QueryInt(q.Get("per_page"), 0))
	response.OK(w, "tags retrieved", result)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	response.OK(w, "tag retrieved", h.svc.Get(r.Context(), chi.URLParam(r, "id")))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in Input
	if err := httpx.Decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}
	response.Created(w, "tag created", h.svc.Create(r.Context(), in))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var in Input
	if err := httpx.Decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}
	response.OK(w, "tag updated", h.svc.Update(r.Context(), chi.URLParam(r, "id"), in))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	h.svc.Delete(r.Context(), chi.URLParam(r, "id"))
	response.OK(w, "tag deleted", nil)
}
