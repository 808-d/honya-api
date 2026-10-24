package middleware

import (
	"net/http"

	"github.com/808-d/honya-api/internal/database/sqlc"
	"github.com/808-d/honya-api/internal/shared/apperror"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RequireFeature allows the request only when one of the caller's roles grants
// the feature code. It must run after Auth so the user id is in the context.
// Roles and their features are read from the database on every request, so
// permission changes take effect immediately.
func RequireFeature(db *pgxpool.Pool, code string) func(http.Handler) http.Handler {
	q := sqlc.New()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := UserIDFromContext(r.Context())
			if userID == "" {
				panic(apperror.Unauthorized("authentication required"))
			}

			has, err := q.UserHasFeature(r.Context(), db, sqlc.UserHasFeatureParams{
				PublicID: userID,
				Code:     code,
			})
			if err != nil {
				panic(err)
			}
			if !has {
				panic(apperror.Forbidden("missing required feature: " + code))
			}

			next.ServeHTTP(w, r)
		})
	}
}
