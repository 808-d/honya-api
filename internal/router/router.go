package router

import (
	"net/http"
	"time"

	"github.com/808-d/honya-api/internal/config"
	appmw "github.com/808-d/honya-api/internal/middleware"
	"github.com/808-d/honya-api/internal/modules/auth"
	"github.com/808-d/honya-api/internal/modules/catalog"
	"github.com/808-d/honya-api/pkg/jwt"
	"github.com/808-d/honya-api/pkg/mailer"
	"github.com/808-d/honya-api/pkg/verification"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// New wires up the HTTP router with its dependencies.
func New(db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.ClientIPFromRemoteAddr)
	r.Use(chimw.Logger)
	r.Use(appmw.ExceptionHandler)
	r.Use(appmw.RateLimit(100, time.Minute))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"message":"ok"}`))
	})

	tokenManager := jwt.NewManager(
		cfg.JWTSecret,
		cfg.AppName,
		time.Duration(cfg.JWTAccessTTLMin)*time.Minute,
	)
	codes := verification.NewStore(rdb, time.Duration(cfg.VerificationTTLMin)*time.Minute)

	authHandler := auth.NewHandler(
		auth.NewService(
			auth.NewRepository(db),
			tokenManager,
			codes,
			mailer.NewDummy(),
			cfg.AppBaseURL,
			cfg.WebBaseURL,
		),
		cfg.AppEnv == "production",
	)

	catalogHandler := catalog.NewHandler(catalog.NewService(catalog.NewRepository(db)))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			authHandler.RegisterRoutes(r)
		})
		r.Route("/catalog", func(r chi.Router) {
			catalogHandler.RegisterRoutes(r)
		})
	})

	return r
}
