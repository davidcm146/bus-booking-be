package configs

import (
	"log/slog"
	"os"
	"time"
)

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

func LoadJWTConfig() JWTConfig {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		if os.Getenv("APP_ENV") == "production" {
			slog.Error("JWT_SECRET is required in production")
		}
	}

	return JWTConfig{
		Secret:    secret,
		ExpiresIn: GetDuration("JWT_EXPIRES_IN", 24*time.Hour),
	}
}
