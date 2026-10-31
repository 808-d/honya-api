package middleware

import (
	"net/http"
	"slices"

	chicors "github.com/go-chi/cors"
)

// CORS lets the web frontend call the API directly from the browser.
//
// When dev is true every Origin is reflected (convenient for localhost on any
// port); otherwise only allowedOrigin is accepted. Credentials are allowed so
// the httpOnly access_token cookie is stored and sent.
func CORS(allowedOrigins []string, dev bool) func(http.Handler) http.Handler {
	return chicors.Handler(chicors.Options{
		AllowOriginFunc: func(_ *http.Request, origin string) bool {
			return origin != "" && (dev || slices.Contains(allowedOrigins, origin))
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           600,
	})
}
