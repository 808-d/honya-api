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

	// URLs used to build verification / reset links
	AppBaseURL string
	WebBaseURL string

	// Verification tokens
	VerificationTTLMin int

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

		AppBaseURL: getEnv("APP_BASE_URL", "http://localhost:8080"),
		WebBaseURL: getEnv("WEB_BASE_URL", "http://localhost:3000"),

		VerificationTTLMin: getInt("VERIFICATION_TTL_MINUTES", 15),

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
