package service

import (
	"context"
	"testing"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/mocks"
	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	t.Run("success with default role", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewUserService(repo)

		repo.On("FindByEmail", mock.Anything, "bob@example.com").Return(nil, apperror.NotFound("user.not_found"))
		repo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Email == "bob@example.com" && u.Role == model.RolePassenger
		})).Return(nil)

		resp, err := svc.CreateUser(context.Background(), request.CreateUserRequest{
			Name:     "Bob",
			Email:    "bob@example.com",
			Password: "password123",
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "Bob", resp.Name)
		require.Equal(t, "bob@example.com", resp.Email)
		require.Equal(t, model.RolePassenger, resp.Role)
	})

	t.Run("email conflict", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewUserService(repo)

		existing := &model.User{ID: 1, Email: "bob@example.com"}
		repo.On("FindByEmail", mock.Anything, "bob@example.com").Return(existing, nil)

		resp, err := svc.CreateUser(context.Background(), request.CreateUserRequest{
			Name:     "Bob",
			Email:    "bob@example.com",
			Password: "password123",
		})

		require.Nil(t, resp)
		require.Error(t, err)
		appErr, ok := err.(*apperror.AppError)
		require.True(t, ok)
		require.Equal(t, 409, appErr.HTTPStatus)
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("custom role", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewUserService(repo)

		repo.On("FindByEmail", mock.Anything, "op@example.com").Return(nil, apperror.NotFound("user.not_found"))
		repo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Role == model.RoleOperator
		})).Return(nil)

		resp, err := svc.CreateUser(context.Background(), request.CreateUserRequest{
			Name:     "Op",
			Email:    "op@example.com",
			Password: "password123",
			Role:     "operator",
		})

		require.NoError(t, err)
		require.Equal(t, model.RoleOperator, resp.Role)
	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewUserService(repo)

		user := &model.User{ID: 3, FullName: "Alice", Email: "alice@example.com"}
		repo.On("FindByID", mock.Anything, 3).Return(user, nil)

		resp, err := svc.GetUserByID(context.Background(), 3)

		require.NoError(t, err)
		require.Equal(t, 3, resp.ID)
		require.Equal(t, "Alice", resp.Name)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewUserService(repo)

		repo.On("FindByID", mock.Anything, 999).Return(nil, apperror.NotFound("user.not_found"))

		resp, err := svc.GetUserByID(context.Background(), 999)

		require.Nil(t, resp)
		require.Error(t, err)
	})
}

func TestGetAllUsers(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	svc := NewUserService(repo)

	users := []model.User{
		{ID: 1, FullName: "A", Email: "a@example.com"},
		{ID: 2, FullName: "B", Email: "b@example.com"},
	}
	repo.On("FindAll", mock.Anything).Return(users, nil)

	resp, err := svc.GetAllUsers(context.Background())

	require.NoError(t, err)
	require.Len(t, resp, 2)
	require.Equal(t, "A", resp[0].Name)
}

func TestUpdateUser(t *testing.T) {
	t.Run("name change only", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewUserService(repo)

		user := &model.User{ID: 1, FullName: "Old", Email: "old@example.com"}
		repo.On("FindByID", mock.Anything, 1).Return(user, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

		resp, err := svc.UpdateUser(context.Background(), 1, request.UpdateUserRequest{
			Name: "New",
		})

		require.NoError(t, err)
		require.Equal(t, "New", resp.Name)
		require.Equal(t, "old@example.com", resp.Email)
	})

	t.Run("email change with no conflict", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewUserService(repo)

		user := &model.User{ID: 1, FullName: "Old", Email: "old@example.com"}
		repo.On("FindByID", mock.Anything, 1).Return(user, nil)
		repo.On("FindByEmail", mock.Anything, "new@example.com").Return(nil, apperror.NotFound("user.not_found"))
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

		resp, err := svc.UpdateUser(context.Background(), 1, request.UpdateUserRequest{
			Email: "new@example.com",
		})

		require.NoError(t, err)
		require.Equal(t, "new@example.com", resp.Email)
	})

	t.Run("email change conflicts with another user", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewUserService(repo)

		user := &model.User{ID: 1, FullName: "Old", Email: "old@example.com"}
		other := &model.User{ID: 2, Email: "new@example.com"}
		repo.On("FindByID", mock.Anything, 1).Return(user, nil)
		repo.On("FindByEmail", mock.Anything, "new@example.com").Return(other, nil)

		resp, err := svc.UpdateUser(context.Background(), 1, request.UpdateUserRequest{
			Email: "new@example.com",
		})

		require.Nil(t, resp)
		require.Error(t, err)
		appErr, ok := err.(*apperror.AppError)
		require.True(t, ok)
		require.Equal(t, 409, appErr.HTTPStatus)
		repo.AssertNotCalled(t, "Update")
	})

	t.Run("email change to own email is not a conflict", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewUserService(repo)

		user := &model.User{ID: 1, FullName: "Old", Email: "old@example.com"}
		repo.On("FindByID", mock.Anything, 1).Return(user, nil)
		repo.On("FindByEmail", mock.Anything, "old@example.com").Return(user, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

		resp, err := svc.UpdateUser(context.Background(), 1, request.UpdateUserRequest{
			Email: "old@example.com",
		})

		require.NoError(t, err)
		require.Equal(t, "old@example.com", resp.Email)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewUserService(repo)

		repo.On("Delete", mock.Anything, 1).Return(nil)

		err := svc.DeleteUser(context.Background(), 1)

		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewUserService(repo)

		repo.On("Delete", mock.Anything, 999).Return(apperror.NotFound("user.not_found"))

		err := svc.DeleteUser(context.Background(), 999)

		require.Error(t, err)
	})
}
