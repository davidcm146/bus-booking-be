package module

import (
	"github.com/davidcm146/bus-booking-be/configs"
	"github.com/davidcm146/bus-booking-be/internal/handler"
	"github.com/davidcm146/bus-booking-be/internal/repository"
	"github.com/davidcm146/bus-booking-be/internal/service"
	"gorm.io/gorm"
)

type AuthModule struct {
	Handler *handler.AuthHandler
}

func NewAuthModule(db *gorm.DB, jwtCfg configs.JWTConfig) *AuthModule {
	svc := service.NewAuthService(repository.NewUserRepository(db), repository.NewRoleRepository(db), jwtCfg)
	return &AuthModule{
		Handler: handler.NewAuthHandler(svc),
	}
}
