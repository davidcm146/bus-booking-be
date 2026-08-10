package mocks

import (
	"context"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/dto/response"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a testify mock for service.AuthService.
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Signup(ctx context.Context, req request.SignupRequest) (*response.SignupResponse, error) {
	ret := m.Called(ctx, req)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*response.SignupResponse), ret.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error) {
	ret := m.Called(ctx, req)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*response.LoginResponse), ret.Error(1)
}

func (m *MockAuthService) Me(ctx context.Context) (*response.UserResponse, error) {
	ret := m.Called(ctx)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*response.UserResponse), ret.Error(1)
}
