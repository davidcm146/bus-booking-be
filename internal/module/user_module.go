package module

import (
	"github.com/davidcm146/bus-booking-be/internal/handler"
	"github.com/davidcm146/bus-booking-be/internal/repository"
	"github.com/davidcm146/bus-booking-be/internal/service"
	"gorm.io/gorm"
)

type UserModule struct {
	Handler *handler.UserHandler
}

func NewUserModule(db *gorm.DB) *UserModule {
	svc := service.NewUserService(repository.NewUserRepository(db), repository.NewRoleRepository(db))
	return &UserModule{
		Handler: handler.NewUserHandler(svc),
	}
}
