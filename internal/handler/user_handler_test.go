package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/dto/response"
	"github.com/davidcm146/bus-booking-be/internal/mocks"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newUserRouter(h *UserHandler) *gin.Engine {
	r := gin.New()
	v1 := r.Group("/api/v1")
	users := v1.Group("/users")
	users.POST("", h.Create)
	users.GET("", h.GetAll)
	users.GET("/:id", h.GetByID)
	users.PUT("/:id", h.Update)
	users.DELETE("/:id", h.Delete)
	return r
}

func TestUserHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockUserService)
		h := NewUserHandler(svc)
		r := newUserRouter(h)

		expected := &response.UserResponse{ID: 1, Name: "Bob", Email: "bob@example.com", Role: "passenger"}
		svc.On("CreateUser", mock.Anything, request.CreateUserRequest{
			Name:     "Bob",
			Email:    "bob@example.com",
			Password: "password123",
		}).Return(expected, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", jsonBody(t, request.CreateUserRequest{
			Name:     "Bob",
			Email:    "bob@example.com",
			Password: "password123",
		}))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("validation fails on empty body", func(t *testing.T) {
		svc := new(mocks.MockUserService)
		h := NewUserHandler(svc)
		r := newUserRouter(h)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		svc.AssertNotCalled(t, "CreateUser")
	})

	t.Run("conflict returns 409", func(t *testing.T) {
		svc := new(mocks.MockUserService)
		h := NewUserHandler(svc)
		r := newUserRouter(h)

		svc.On("CreateUser", mock.Anything, mock.Anything).Return(nil, apperror.Conflict("user.email_conflict"))

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", jsonBody(t, request.CreateUserRequest{
			Name:     "Bob",
			Email:    "bob@example.com",
			Password: "password123",
		}))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusConflict, rec.Code)
	})
}

func TestUserHandler_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockUserService)
		h := NewUserHandler(svc)
		r := newUserRouter(h)

		expected := &response.UserResponse{ID: 3, Name: "Alice", Email: "alice@example.com"}
		svc.On("GetUserByID", mock.Anything, 3).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/3", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := new(mocks.MockUserService)
		h := NewUserHandler(svc)
		r := newUserRouter(h)

		svc.On("GetUserByID", mock.Anything, 999).Return(nil, apperror.NotFound("user.not_found"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/999", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid id returns 422", func(t *testing.T) {
		svc := new(mocks.MockUserService)
		h := NewUserHandler(svc)
		r := newUserRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/abc", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		svc.AssertNotCalled(t, "GetUserByID")
	})
}

func TestUserHandler_GetAll(t *testing.T) {
	svc := new(mocks.MockUserService)
	h := NewUserHandler(svc)
	r := newUserRouter(h)

	expected := []response.UserResponse{
		{ID: 1, Name: "A", Email: "a@example.com"},
		{ID: 2, Name: "B", Email: "b@example.com"},
	}
	svc.On("GetAllUsers", mock.Anything).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUserHandler_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockUserService)
		h := NewUserHandler(svc)
		r := newUserRouter(h)

		expected := &response.UserResponse{ID: 1, Name: "New", Email: "old@example.com"}
		svc.On("UpdateUser", mock.Anything, 1, request.UpdateUserRequest{
			Name: "New",
		}).Return(expected, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/1", jsonBody(t, request.UpdateUserRequest{
			Name: "New",
		}))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := new(mocks.MockUserService)
		h := NewUserHandler(svc)
		r := newUserRouter(h)

		svc.On("UpdateUser", mock.Anything, 999, mock.Anything).Return(nil, apperror.NotFound("user.not_found"))

		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/999", jsonBody(t, request.UpdateUserRequest{
			Name: "New",
		}))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestUserHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockUserService)
		h := NewUserHandler(svc)
		r := newUserRouter(h)

		svc.On("DeleteUser", mock.Anything, 1).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := new(mocks.MockUserService)
		h := NewUserHandler(svc)
		r := newUserRouter(h)

		svc.On("DeleteUser", mock.Anything, 999).Return(apperror.NotFound("user.not_found"))

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/999", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid id returns 422", func(t *testing.T) {
		svc := new(mocks.MockUserService)
		h := NewUserHandler(svc)
		r := newUserRouter(h)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/abc", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		svc.AssertNotCalled(t, "DeleteUser")
	})
}
