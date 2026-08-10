package service

import (
	"context"
	"testing"

	"github.com/davidcm146/bus-booking-be/configs"
	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/mocks"
	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/auth"
	"github.com/davidcm146/bus-booking-be/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Dummy values for tests only — not real credentials.
// Must satisfy the strong_password validator (upper, lower, digit, 8+ chars).
const testPassword = "Testpass1dummy"
const wrongPassword = "Wrongpass1dummy"

func newAuthSvc(repo *mocks.MockUserRepository) *AuthServiceImpl {
	return NewAuthService(repo, configs.JWTConfig{
		Secret:    "test-secret",
		ExpiresIn: 0,
	})
}

func TestSignup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthSvc(repo)

		repo.On("FindByPhone", mock.Anything, "+84901234567").Return(nil, apperror.NotFound("user.not_found"))
		repo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

		resp, err := svc.Signup(context.Background(), request.SignupRequest{
			Phone:    "+84901234567",
			Password: testPassword,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, resp.Token)
		require.Equal(t, "+84901234567", resp.Phone)
		repo.AssertExpectations(t)
	})

	t.Run("phone conflict", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthSvc(repo)

		existing := &model.User{ID: 1, Phone: "+84901234567"}
		repo.On("FindByPhone", mock.Anything, "+84901234567").Return(existing, nil)

		resp, err := svc.Signup(context.Background(), request.SignupRequest{
			Phone:    "+84901234567",
			Password: testPassword,
		})

		require.Nil(t, resp)
		require.Error(t, err)
		appErr, ok := err.(*apperror.AppError)
		require.True(t, ok)
		require.Equal(t, 409, appErr.HTTPStatus)
		repo.AssertNotCalled(t, "Create")
	})
}

func TestLogin(t *testing.T) {
	hashed, err := utils.HashPassword(testPassword)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := newAuthSvc(repo)

		user := &model.User{ID: 5, Phone: "+84901234567", Password: hashed}
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

		repo.On("FindByPhone", mock.Anything, "+84900000000").Return(nil, apperror.NotFound("user.not_found"))

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

		user := &model.User{ID: 5, Phone: "+84901234567", Password: hashed}
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

		user := &model.User{ID: 7, Phone: "+84901234567"}
		repo.On("FindByID", mock.Anything, 7).Return(user, nil)

		ctx := auth.WithContext(context.Background(), &auth.Context{UserID: 7, Role: "passenger"})

		resp, err := svc.Me(ctx)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, 7, resp.ID)
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
