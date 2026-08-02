package configs

type AppConfig struct {
	Environment string
	Port        int
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Environment: GetString("APP_ENV", "development"),
		Port:        GetInt("APP_PORT", 8080),
	}
}
