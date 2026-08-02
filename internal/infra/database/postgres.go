package database

import (
	"database/sql"

	"github.com/davidcm146/bus-booking-be/configs"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg configs.DatabaseConfig) (*sql.DB, error) {

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
