package middleware

import (
	"net/http"
	"time"

	"github.com/808-d/honya-api/internal/shared/response"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

// RateLimit throttles requests per client IP using chi's httprate middleware,
// rendering exceeded limits through the standard JSON envelope.
func RateLimit(requests int, window time.Duration) func(http.Handler) http.Handler {
	return httprate.LimitBy(
		requests,
		window,
		func(r *http.Request) (string, error) {
			return httprate.CanonicalizeIP(chimw.GetClientIP(r.Context())), nil
		},
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			response.Error(w, http.StatusTooManyRequests, "too_many_requests", "rate limit exceeded")
		}),
	)
}
