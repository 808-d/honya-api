package order

import (
	"net/http"

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

// Checkout handles POST /orders.
func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	var in CreateOrderRequest
	if err := httpx.Decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	result := h.svc.Checkout(r.Context(), middleware.UserIDFromContext(r.Context()), in)
	response.Created(w, "order created", result)
}

// List handles GET /orders.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	result := h.svc.ListForUser(
		r.Context(),
		middleware.UserIDFromContext(r.Context()),
		ListFilter{
			Status:  httpx.QueryInt16(query.Get("status")),
			Page:    httpx.QueryInt(query.Get("page"), 1),
			PerPage: httpx.QueryInt(query.Get("per_page"), 0),
		},
	)
	response.OK(w, "orders retrieved", result)
}

// Get handles GET /orders/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	result := h.svc.GetForUser(
		r.Context(),
		middleware.UserIDFromContext(r.Context()),
		chi.URLParam(r, "id"),
	)
	response.OK(w, "order retrieved", result)
}

// Cancel handles POST /orders/{id}/cancel.
func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	result := h.svc.Cancel(
		r.Context(),
		middleware.UserIDFromContext(r.Context()),
		chi.URLParam(r, "id"),
	)
	response.OK(w, "order cancelled", result)
}

// AdminList handles GET /admin/orders.
func (h *Handler) AdminList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	result := h.svc.ListAdmin(r.Context(), ListFilter{
		Status:        httpx.QueryInt16(query.Get("status")),
		PaymentStatus: httpx.QueryInt16(query.Get("payment_status")),
		Page:          httpx.QueryInt(query.Get("page"), 1),
		PerPage:       httpx.QueryInt(query.Get("per_page"), 0),
	})
	response.OK(w, "orders retrieved", result)
}

// AdminGet handles GET /admin/orders/{id}.
func (h *Handler) AdminGet(w http.ResponseWriter, r *http.Request) {
	result := h.svc.GetAdmin(r.Context(), chi.URLParam(r, "id"))
	response.OK(w, "order retrieved", result)
}

// AdminUpdateStatus handles PATCH /admin/orders/{id}/status.
func (h *Handler) AdminUpdateStatus(w http.ResponseWriter, r *http.Request) {
	var in UpdateStatusRequest
	if err := httpx.Decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	result := h.svc.UpdateStatus(r.Context(), chi.URLParam(r, "id"), in)
	response.OK(w, "order status updated", result)
}

// AdminUpdatePayment handles PATCH /admin/orders/{id}/payment.
func (h *Handler) AdminUpdatePayment(w http.ResponseWriter, r *http.Request) {
	var in UpdatePaymentRequest
	if err := httpx.Decode(r, &in); err != nil {
		panic(apperror.BadRequest("invalid request body"))
	}

	result := h.svc.UpdatePayment(r.Context(), chi.URLParam(r, "id"), in)
	response.OK(w, "order payment updated", result)
}
