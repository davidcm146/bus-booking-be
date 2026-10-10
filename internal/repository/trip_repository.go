package repository

import (
	"context"
	"time"

	"github.com/davidcm146/bus-booking-be/internal/model"
)

type TripSearchFilter struct {
	OriginAreaID      int
	DestinationAreaID int
	DepartureFrom     time.Time
	DepartureTo       time.Time
	ArrivalFrom       time.Time
	ArrivalTo         time.Time
	Limit             int
	Offset            int
}

type TripRepository interface {
	Search(ctx context.Context, filter TripSearchFilter) ([]model.Trip, error)
	FindByID(ctx context.Context, id int) (*model.Trip, error)
}
