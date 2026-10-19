package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/808-d/honya-api/internal/shared/apperror"
	"github.com/808-d/honya-api/pkg/jwt"
)

const accessTokenCookie = "access_token"

type contextKey int

const userIDKey contextKey = iota

// Auth requires a valid access token (httpOnly cookie or Bearer header) and
// stores the caller's public user id in the request context. It panics with an
// *apperror.Error so the global exception handler renders a 401.
func Auth(tokens *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := accessToken(r)
			if raw == "" {
				panic(apperror.Unauthorized("authentication required"))
			}

			claims, err := tokens.Parse(raw)
			if err != nil {
				panic(apperror.Unauthorized("invalid or expired token"))
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext returns the authenticated user's public id, or "" when the
// request is not authenticated.
func UserIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(userIDKey).(string)
	return id
}

func accessToken(r *http.Request) string {
	if c, err := r.Cookie(accessTokenCookie); err == nil && c.Value != "" {
		return c.Value
	}
	if header := r.Header.Get("Authorization"); strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ""
}
