package middleware

import (
	"log"
	"net/http"

	"github.com/808-d/honya-api/internal/shared/apperror"
	"github.com/808-d/honya-api/internal/shared/response"
)

// ExceptionHandler recovers from panics raised anywhere in the handler chain.
// An *apperror.Error is rendered with its status/code; any other panic becomes
// a 500 and is logged.
func ExceptionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}

			if appErr, ok := rec.(*apperror.Error); ok {
				response.Error(w, appErr.Status, appErr.Code, appErr.Message)
				return
			}

			log.Printf("panic recovered: %v", rec)
			response.Internal(w)
		}()

		next.ServeHTTP(w, r)
	})
}
