package module

import (
	"github.com/davidcm146/bus-booking-be/internal/handler"
	"github.com/davidcm146/bus-booking-be/internal/repository"
	"github.com/davidcm146/bus-booking-be/internal/service"
	"gorm.io/gorm"
)

type TripModule struct {
	Handler *handler.TripHandler
}

func NewTripModule(db *gorm.DB) *TripModule {
	svc := service.NewTripService(repository.NewTripRepository(db))
	return &TripModule{
		Handler: handler.NewTripHandler(svc),
	}
}
