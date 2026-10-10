package service

import (
	"context"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/dto/response"
)

// TripService defines the business-logic contract for Trip queries.
type TripService interface {
	SearchTrips(ctx context.Context, req request.SearchTripsRequest) ([]response.TripResponse, error)
	GetTripByID(ctx context.Context, id int) (*response.TripResponse, error)
}
