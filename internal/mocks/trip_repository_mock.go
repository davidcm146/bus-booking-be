package mocks

import (
	"context"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/repository"
	"github.com/stretchr/testify/mock"
)

// MockTripRepository is a testify mock for repository.TripRepository.
type MockTripRepository struct {
	mock.Mock
}

func (m *MockTripRepository) Search(ctx context.Context, filter repository.TripSearchFilter) ([]model.Trip, error) {
	ret := m.Called(ctx, filter)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]model.Trip), ret.Error(1)
}

func (m *MockTripRepository) FindByID(ctx context.Context, id int) (*model.Trip, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*model.Trip), ret.Error(1)
}
