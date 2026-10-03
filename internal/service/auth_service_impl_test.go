package service

import (
	"context"
	"testing"
	"time"

	"github.com/davidcm146/bus-booking-be/configs"
	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/mocks"
	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/auth"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/davidcm146/bus-booking-be/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Dummy values for tests only — not real credentials.
// Must satisfy the strong_password validator (upper, lower, digit, 8+ chars).
const testPassword = "Testpass1dummy"
const wrongPassword = "Wrongpass1dummy"

func newAuthSvc(repo *mocks.MockUserRepository) *AuthServiceImpl {
	roles := new(mocks.MockRoleRepository)
	roles.On("FindByName", mock.Anything, model.RolePassenger).Return(passengerRole, nil)
	return NewAuthService(repo, roles, configs.JWTConfig{
		Secret:    "test-secret",
		ExpiresIn: time.Hour,
	})
}

func TestSignup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthSvc(repo)

		repo.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return *u.Phone == "+84901234567" && *u.RoleID == passengerRole.ID && u.Status == model.UserStatusActive
		})).Return(nil)

		resp, err := svc.Signup(context.Background(), request.SignupRequest{
			Phone:    "+84901234567",
			Password: testPassword,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, resp.Token)
		require.Equal(t, "+84901234567", resp.Phone)
		require.Equal(t, model.RolePassenger, resp.Role)
		claims, err := utils.ParseToken(resp.Token, "test-secret")
		require.NoError(t, err)
		require.Equal(t, string(model.RolePassenger), claims.Role)
		repo.AssertExpectations(t)
	})

	t.Run("phone conflict", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthSvc(repo)

		repo.On("Create", mock.Anything, mock.Anything).Return(apperror.Conflict(i18n.KeyUserPhoneConflict))

		resp, err := svc.Signup(context.Background(), request.SignupRequest{
			Phone:    "+84901234567",
			Password: testPassword,
		})

		require.Nil(t, resp)
		require.Error(t, err)
		appErr, ok := err.(*apperror.AppError)
		require.True(t, ok)
		require.Equal(t, 409, appErr.HTTPStatus)
	})
}

func TestLogin(t *testing.T) {
	hashed, err := utils.HashPassword(testPassword)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthSvc(repo)

		user := &model.User{ID: 5, Phone: new("+84901234567"), Password: hashed, Role: passengerRole}
		repo.On("FindByPhone", mock.Anything, "+84901234567").Return(user, nil)

		resp, err := svc.Login(context.Background(), request.LoginRequest{
			Phone:    "+84901234567",
			Password: testPassword,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, resp.Token)
		require.Equal(t, 5, resp.ID)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthSvc(repo)

		repo.On("FindByPhone", mock.Anything, "+84900000000").Return(nil, apperror.NotFound(i18n.KeyUserNotFound))

		resp, err := svc.Login(context.Background(), request.LoginRequest{
			Phone:    "+84900000000",
			Password: testPassword,
		})

		require.Nil(t, resp)
		require.Error(t, err)
		appErr, ok := err.(*apperror.AppError)
		require.True(t, ok)
		require.Equal(t, 401, appErr.HTTPStatus)
	})

	t.Run("wrong password", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthSvc(repo)

		user := &model.User{ID: 5, Phone: new("+84901234567"), Password: hashed, Role: passengerRole}
		repo.On("FindByPhone", mock.Anything, "+84901234567").Return(user, nil)

		resp, err := svc.Login(context.Background(), request.LoginRequest{
			Phone:    "+84901234567",
			Password: wrongPassword,
		})

		require.Nil(t, resp)
		require.Error(t, err)
		appErr, ok := err.(*apperror.AppError)
		require.True(t, ok)
		require.Equal(t, 401, appErr.HTTPStatus)
	})
}

func TestMe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthSvc(repo)

		user := &model.User{ID: 7, Phone: new("+84901234567")}
		repo.On("FindByID", mock.Anything, 7).Return(user, nil)

		ctx := auth.WithContext(context.Background(), &auth.Context{UserID: 7, Role: "passenger"})

		resp, err := svc.Me(ctx)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, 7, resp.ID)
	})

	t.Run("deleted or missing user is not found", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthSvc(repo)

		repo.On("FindByID", mock.Anything, 7).Return(nil, apperror.NotFound(i18n.KeyUserNotFound))

		ctx := auth.WithContext(context.Background(), &auth.Context{UserID: 7, Role: "passenger"})

		resp, err := svc.Me(ctx)

		require.Nil(t, resp)
		appErr, ok := err.(*apperror.AppError)
		require.True(t, ok)
		require.Equal(t, 404, appErr.HTTPStatus)
	})

	t.Run("no auth context", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthSvc(repo)

		resp, err := svc.Me(context.Background())

		require.Nil(t, resp)
		require.Error(t, err)
		appErr, ok := err.(*apperror.AppError)
		require.True(t, ok)
		require.Equal(t, 401, appErr.HTTPStatus)
	})
}
