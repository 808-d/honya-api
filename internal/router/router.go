package router

import (
	"net/http"
	"time"

	"github.com/808-d/honya-api/internal/config"
	appmw "github.com/808-d/honya-api/internal/middleware"
	"github.com/808-d/honya-api/internal/modules/auth"
	"github.com/808-d/honya-api/internal/modules/author"
	"github.com/808-d/honya-api/internal/modules/cart"
	"github.com/808-d/honya-api/internal/modules/coupon"
	"github.com/808-d/honya-api/internal/modules/dashboard"
	"github.com/808-d/honya-api/internal/modules/genre"
	"github.com/808-d/honya-api/internal/modules/order"
	"github.com/808-d/honya-api/internal/modules/product"
	"github.com/808-d/honya-api/internal/modules/publisher"
	"github.com/808-d/honya-api/internal/modules/rbac"
	"github.com/808-d/honya-api/internal/modules/tag"
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
	r.Use(appmw.CORS([]string{cfg.WebBaseURL, cfg.DeployURL}, cfg.AppEnv != "production"))
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
	transactionalMail := mailer.New(cfg.SMTPEmail, cfg.AppPassword)

	authHandler := auth.NewHandler(
		auth.NewService(
			auth.NewRepository(db),
			tokenManager,
			codes,
			transactionalMail,
			cfg.AppBaseURL,
			cfg.WebBaseURL,
		),
		cfg.AppEnv == "production",
	)

	catalogWrite := appmw.RequireFeature(db, "catalog.write")
	productHandler := product.NewHandler(product.NewService(product.NewRepository(db)))
	authorHandler := author.NewHandler(author.NewService(author.NewRepository(db)))
	genreHandler := genre.NewHandler(genre.NewService(genre.NewRepository(db)))
	publisherHandler := publisher.NewHandler(publisher.NewService(publisher.NewRepository(db)))
	tagHandler := tag.NewHandler(tag.NewService(tag.NewRepository(db)))
	cartHandler := cart.NewHandler(cart.NewService(cart.NewRepository(db)))
	couponService := coupon.NewService(coupon.NewRepository(db))
	couponHandler := coupon.NewHandler(couponService)
	orderHandler := order.NewHandler(order.NewService(
		order.NewRepository(db),
		couponService,
		cfg.ShippingFeePercent,
		transactionalMail,
		cfg.WebBaseURL,
	))
	rbacHandler := rbac.NewHandler(rbac.NewService(rbac.NewRepository(db)))
	dashboardHandler := dashboard.NewHandler(dashboard.NewService(dashboard.NewRepository(db), rdb))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			authHandler.RegisterRoutes(r)
		})
		r.With(appmw.Auth(tokenManager)).Get("/auth/permissions", authHandler.Permissions)
		authMiddleware := appmw.Auth(tokenManager)
		productHandler.RegisterRoutes(r, authMiddleware, catalogWrite)
		authorHandler.RegisterRoutes(r, authMiddleware, catalogWrite)
		genreHandler.RegisterRoutes(r, authMiddleware, catalogWrite)
		publisherHandler.RegisterRoutes(r, authMiddleware, catalogWrite)
		tagHandler.RegisterRoutes(r, authMiddleware, catalogWrite)
		couponHandler.RegisterRoutes(
			r,
			appmw.Auth(tokenManager),
			appmw.RequireFeature(db, "coupon.read"),
			appmw.RequireFeature(db, "coupon.write"),
		)
		orderHandler.RegisterRoutes(
			r,
			appmw.Auth(tokenManager),
			appmw.RequireFeature(db, "order.read"),
			appmw.RequireFeature(db, "order.write"),
		)
		rbacHandler.RegisterRoutes(
			r,
			appmw.Auth(tokenManager),
			appmw.RequireFeature(db, "role.read"),
			appmw.RequireFeature(db, "role.write"),
			appmw.RequireFeature(db, "user.read"),
			appmw.RequireFeature(db, "user.write"),
		)
		dashboardHandler.RegisterRoutes(
			r,
			appmw.Auth(tokenManager),
			appmw.RequireFeature(db, "dashboard.access"),
		)

		r.Group(func(r chi.Router) {
			r.Use(appmw.Auth(tokenManager))
			cartHandler.RegisterRoutes(r)
		})
	})

	return r
}
