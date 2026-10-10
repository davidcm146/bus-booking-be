package configs

import "strings"

type AppConfig struct {
	Environment string
	Port        int
	CORSOrigins []string
}

func LoadAppConfig() AppConfig {
	origins := GetString("CORS_ORIGINS", "http://localhost:3000,http://localhost:3001")
	return AppConfig{
		Environment: GetString("APP_ENV", "development"),
		Port:        GetInt("APP_PORT", 8080),
		CORSOrigins: strings.Split(origins, ","),
	}
}
