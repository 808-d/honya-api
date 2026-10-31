package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppEnv  string
	Port    string

	DatabaseURL string

	// JWT
	JWTSecret       string
	JWTAccessTTLMin int

	// Frontend URLs used to build the verification / reset links in emails
	AppBaseURL string
	WebBaseURL string
	DeployURL  string

	// Verification tokens
	VerificationTTLMin int

	// Percentage of the pre-discount subtotal charged as shipping.
	ShippingFeePercent float64

	// Gmail transactional mail (leave empty to log instead of sending)
	SMTPEmail   string
	AppPassword string

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		AppName: getEnv("APP_NAME", "honya-api"),
		AppEnv:  getEnv("APP_ENV", "development"),
		Port:    getEnv("PORT", "8080"),

		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:123@localhost:5432/honya?sslmode=disable"),

		JWTSecret:       getEnv("JWT_SECRET", "change-me-super-secret"),
		JWTAccessTTLMin: getInt("JWT_ACCESS_TTL_MINUTES", 15),

		AppBaseURL: getEnv("APP_BASE_URL", "http://localhost:5173"),
		WebBaseURL: getEnv("WEB_BASE_URL", "http://localhost:5173"),
		DeployURL:  getEnv("DEPLOY_URL", "http://localhost:5173"),

		VerificationTTLMin: getInt("VERIFICATION_TTL_MINUTES", 30),

		ShippingFeePercent: getFloat("SHIPPING_FEE_PERCENT", 10),

		SMTPEmail:   getEnv("SMTP_EMAIL", ""),
		AppPassword: getEnv("APP_PASSWORD", ""),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getInt("REDIS_DB", 0),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			return n
		}
	}
	return fallback
}
