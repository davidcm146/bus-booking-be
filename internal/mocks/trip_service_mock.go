package mocks

import (
	"context"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/dto/response"
	"github.com/stretchr/testify/mock"
)

// MockTripService is a testify mock for service.TripService.
type MockTripService struct {
	mock.Mock
}

func (m *MockTripService) SearchTrips(ctx context.Context, req request.SearchTripsRequest) ([]response.TripResponse, error) {
	ret := m.Called(ctx, req)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]response.TripResponse), ret.Error(1)
}

func (m *MockTripService) GetTripByID(ctx context.Context, id int) (*response.TripResponse, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*response.TripResponse), ret.Error(1)
}
