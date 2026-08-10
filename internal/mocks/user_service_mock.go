package mocks

import (
	"context"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/dto/response"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a testify mock for service.UserService.
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, req request.CreateUserRequest) (*response.UserResponse, error) {
	ret := m.Called(ctx, req)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*response.UserResponse), ret.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id int) (*response.UserResponse, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*response.UserResponse), ret.Error(1)
}

func (m *MockUserService) GetAllUsers(ctx context.Context) ([]response.UserResponse, error) {
	ret := m.Called(ctx)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]response.UserResponse), ret.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id int, req request.UpdateUserRequest) (*response.UserResponse, error) {
	ret := m.Called(ctx, id, req)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*response.UserResponse), ret.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id int) error {
	ret := m.Called(ctx, id)
	return ret.Error(0)
}
