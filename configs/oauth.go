package configs

import (
	"log/slog"
	"os"
)

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func LoadOAuthConfig() OAuthConfig {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := GetString("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback")

	if clientID == "" && os.Getenv("APP_ENV") == "production" {
		slog.Error("GOOGLE_CLIENT_ID is required in production")
	}

	return OAuthConfig{
		GoogleClientID:     clientID,
		GoogleClientSecret: clientSecret,
		GoogleRedirectURL:  redirectURL,
	}
}
