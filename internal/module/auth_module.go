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

func NewAuthModule(db *gorm.DB, jwtCfg configs.JWTConfig, oauthCfg configs.OAuthConfig) *AuthModule {
	repo := repository.NewUserRepository(db)
	svc := service.NewAuthService(repo, jwtCfg, oauthCfg)
	return &AuthModule{
		Handler: handler.NewAuthHandler(svc, oauthCfg),
	}
}
