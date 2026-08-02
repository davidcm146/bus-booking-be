package configs

import (
	"fmt"
	"time"
)

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	TimeZone        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:            GetString("DB_HOST", "localhost"),
		Port:            GetInt("DB_PORT", 5432),
		User:            GetString("DB_USER", "postgres"),
		Password:        GetString("DB_PASSWORD", "10042003"),
		Name:            GetString("DB_NAME", "bus_booking"),
		SSLMode:         GetString("DB_SSLMODE", "disable"),
		TimeZone:        GetString("DB_TIMEZONE", "Asia/Ho_Chi_Minh"),
		MaxOpenConns:    GetInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    GetInt("DB_MAX_IDLE_CONNS", 10),
		ConnMaxLifetime: GetDuration("DB_CONN_MAX_LIFETIME", time.Hour),
		ConnMaxIdleTime: GetDuration("DB_CONN_MAX_IDLE_TIME", 15*time.Minute),
	}
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
		c.SSLMode,
		c.TimeZone,
	)
}
