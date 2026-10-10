package mocks

import (
	"context"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a testify mock for repository.UserRepository.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	ret := m.Called(ctx, user)
	return ret.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id int) (*model.User, error) {
	ret := m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*model.User), ret.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	ret := m.Called(ctx, email)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*model.User), ret.Error(1)
}

func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	ret := m.Called(ctx, phone)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*model.User), ret.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	ret := m.Called(ctx)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]model.User), ret.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	ret := m.Called(ctx, user)
	return ret.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	ret := m.Called(ctx, id)
	return ret.Error(0)
}
