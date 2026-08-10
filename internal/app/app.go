package app

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/davidcm146/bus-booking-be/configs"
	"github.com/davidcm146/bus-booking-be/internal/infra/database"
	"github.com/davidcm146/bus-booking-be/internal/infra/logger"
	"github.com/davidcm146/bus-booking-be/internal/router"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/davidcm146/bus-booking-be/internal/shared/validator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Application is the root struct that owns all long-lived resources.
type Application struct {
	Config *configs.Config
	DB     *sql.DB
	Gorm   *gorm.DB
	Router *gin.Engine
	Logger *slog.Logger
}

// New bootstraps the entire application: config → DB → DI container → router.
func New() (*Application, error) {
	cfg := configs.Load()
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	log := logger.New(cfg.App.Environment)
	slog.SetDefault(log)

	slog.Info("Connecting to database...")
	i18n.Init()

	sqlDB, err := database.NewPostgres(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("Postgres connection failed: %w", err)
	}

	gormDB, err := database.NewGorm(sqlDB)
	if err != nil {
		return nil, fmt.Errorf("Gorm initialization failed: %w", err)
	}

	slog.Info("Database connected successfully")

	validator.Setup()

	container := NewContainer(gormDB, cfg)

	r := router.New(container.Handlers(), cfg.JWT.Secret, cfg.App.CORSOrigins)

	return &Application{
		Config: cfg,
		DB:     sqlDB,
		Gorm:   gormDB,
		Router: r,
		Logger: log,
	}, nil
}

// Run starts the HTTP server on the configured port.
func (a *Application) Run() error {
	addr := fmt.Sprintf(":%d", a.Config.App.Port)
	slog.Info("Server starting", "addr", addr)
	return a.Router.Run(addr)
}

// Shutdown performs graceful cleanup of resources.
func (a *Application) Shutdown() {
	if a.DB != nil {
		slog.Info("Closing database connection...")
		a.DB.Close()
	}
}
