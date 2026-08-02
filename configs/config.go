package configs

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

func Load() *Config {
	return &Config{
		App:      LoadAppConfig(),
		Database: LoadDatabaseConfig(),
		JWT:      LoadJWTConfig(),
	}
}
