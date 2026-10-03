package service

import (
	"context"
	"testing"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/mocks"
	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/davidcm146/bus-booking-be/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	passengerRole = &model.Role{ID: 1, Name: model.RolePassenger}
	operatorRole  = &model.Role{ID: 2, Name: model.RoleOperator}
)

func newUserSvc() (UserService, *mocks.MockUserRepository, *mocks.MockRoleRepository) {
	userRepo, roleRepo := new(mocks.MockUserRepository), new(mocks.MockRoleRepository)
	return NewUserService(userRepo, roleRepo), userRepo, roleRepo
}

func TestCreateUser(t *testing.T) {
	t.Run("success with default role and hashed password", func(t *testing.T) {
		svc, repo, roles := newUserSvc()

		roles.On("FindByName", mock.Anything, model.RolePassenger).Return(passengerRole, nil)
		repo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return *u.Email == "bob@example.com" && *u.RoleID == passengerRole.ID &&
				u.Status == model.UserStatusActive && u.Phone == nil &&
				u.Password != "password123" && utils.VerifyPassword("password123", u.Password)
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
		repo.AssertExpectations(t)
	})

	t.Run("email conflict from repository is returned", func(t *testing.T) {
		svc, repo, roles := newUserSvc()

		roles.On("FindByName", mock.Anything, model.RolePassenger).Return(passengerRole, nil)
		repo.On("Create", mock.Anything, mock.Anything).Return(apperror.Conflict(i18n.KeyUserEmailConflict))

		resp, err := svc.CreateUser(context.Background(), request.CreateUserRequest{
			Name:     "Bob",
			Email:    "bob@example.com",
			Password: "password123",
		})

		require.Nil(t, resp)
		appErr, ok := err.(*apperror.AppError)
		require.True(t, ok)
		require.Equal(t, 409, appErr.HTTPStatus)
	})

	t.Run("custom role", func(t *testing.T) {
		svc, repo, roles := newUserSvc()

		roles.On("FindByName", mock.Anything, model.RoleOperator).Return(operatorRole, nil)
		repo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return *u.RoleID == operatorRole.ID
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

	t.Run("role lookup failure stops creation", func(t *testing.T) {
		svc, repo, roles := newUserSvc()

		roles.On("FindByName", mock.Anything, model.RolePassenger).Return(nil, apperror.Internal(i18n.KeyRoleFailedFind))

		resp, err := svc.CreateUser(context.Background(), request.CreateUserRequest{
			Name:     "Bob",
			Email:    "bob@example.com",
			Password: "password123",
		})

		require.Nil(t, resp)
		require.Error(t, err)
		repo.AssertNotCalled(t, "Create")
	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, repo, _ := newUserSvc()

		user := &model.User{ID: 3, Name: "Alice", Email: new("alice@example.com")}
		repo.On("FindByID", mock.Anything, 3).Return(user, nil)

		resp, err := svc.GetUserByID(context.Background(), 3)

		require.NoError(t, err)
		require.Equal(t, 3, resp.ID)
		require.Equal(t, "Alice", resp.Name)
	})

	t.Run("not found", func(t *testing.T) {
		svc, repo, _ := newUserSvc()

		repo.On("FindByID", mock.Anything, 999).Return(nil, apperror.NotFound(i18n.KeyUserNotFound))

		resp, err := svc.GetUserByID(context.Background(), 999)

		require.Nil(t, resp)
		require.Error(t, err)
	})
}

func TestGetAllUsers(t *testing.T) {
	svc, repo, _ := newUserSvc()

	users := []model.User{
		{ID: 1, Name: "A", Email: new("a@example.com")},
		{ID: 2, Name: "B", Email: new("b@example.com")},
	}
	repo.On("FindAll", mock.Anything).Return(users, nil)

	resp, err := svc.GetAllUsers(context.Background())

	require.NoError(t, err)
	require.Len(t, resp, 2)
	require.Equal(t, "A", resp[0].Name)
}

func TestUpdateUser(t *testing.T) {
	t.Run("name change only", func(t *testing.T) {
		svc, repo, _ := newUserSvc()

		user := &model.User{ID: 1, Name: "Old", Email: new("old@example.com")}
		repo.On("FindByID", mock.Anything, 1).Return(user, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

		resp, err := svc.UpdateUser(context.Background(), 1, request.UpdateUserRequest{
			Name: "New",
		})

		require.NoError(t, err)
		require.Equal(t, "New", resp.Name)
		require.Equal(t, "old@example.com", resp.Email)
	})

	t.Run("email change", func(t *testing.T) {
		svc, repo, _ := newUserSvc()

		user := &model.User{ID: 1, Name: "Old", Email: new("old@example.com")}
		repo.On("FindByID", mock.Anything, 1).Return(user, nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return *u.Email == "new@example.com"
		})).Return(nil)

		resp, err := svc.UpdateUser(context.Background(), 1, request.UpdateUserRequest{
			Email: "new@example.com",
		})

		require.NoError(t, err)
		require.Equal(t, "new@example.com", resp.Email)
	})

	t.Run("email conflict from repository is returned", func(t *testing.T) {
		svc, repo, _ := newUserSvc()

		user := &model.User{ID: 1, Name: "Old", Email: new("old@example.com")}
		repo.On("FindByID", mock.Anything, 1).Return(user, nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(apperror.Conflict(i18n.KeyUserEmailConflict))

		resp, err := svc.UpdateUser(context.Background(), 1, request.UpdateUserRequest{
			Email: "new@example.com",
		})

		require.Nil(t, resp)
		appErr, ok := err.(*apperror.AppError)
		require.True(t, ok)
		require.Equal(t, 409, appErr.HTTPStatus)
	})

	t.Run("role change resolves role id", func(t *testing.T) {
		svc, repo, roles := newUserSvc()

		user := &model.User{ID: 1, Name: "Old", RoleID: &passengerRole.ID, Role: passengerRole}
		repo.On("FindByID", mock.Anything, 1).Return(user, nil)
		roles.On("FindByName", mock.Anything, model.RoleOperator).Return(operatorRole, nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return *u.RoleID == operatorRole.ID
		})).Return(nil)

		resp, err := svc.UpdateUser(context.Background(), 1, request.UpdateUserRequest{
			Role: "operator",
		})

		require.NoError(t, err)
		require.Equal(t, model.RoleOperator, resp.Role)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, repo, _ := newUserSvc()

		repo.On("Delete", mock.Anything, 1).Return(nil)

		require.NoError(t, svc.DeleteUser(context.Background(), 1))
	})

	t.Run("not found", func(t *testing.T) {
		svc, repo, _ := newUserSvc()

		repo.On("Delete", mock.Anything, 999).Return(apperror.NotFound(i18n.KeyUserNotFound))

		require.Error(t, svc.DeleteUser(context.Background(), 999))
	})
}
