package mocks

import (
	"context"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockRoleRepository is a testify mock for repository.RoleRepository.
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindByName(ctx context.Context, name model.UserRole) (*model.Role, error) {
	ret := m.Called(ctx, name)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*model.Role), ret.Error(1)
}
