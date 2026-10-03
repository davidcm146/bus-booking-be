package configs

import "time"

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

func LoadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:    GetString("JWT_SECRET", ""),
		ExpiresIn: GetDuration("JWT_EXPIRES_IN", 24*time.Hour),
	}
}
