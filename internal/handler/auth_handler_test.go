package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/dto/response"
	"github.com/davidcm146/bus-booking-be/internal/mocks"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/auth"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/davidcm146/bus-booking-be/internal/shared/validator"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestMain initializes i18n and validator once for all handler tests.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	i18nSvc, err := i18n.New()
	if err != nil {
		panic("failed to init i18n for tests: " + err.Error())
	}
	i18n.Setup(i18nSvc)
	validator.Setup()

	m.Run()
}

func jsonBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewBuffer(b)
}

func multipartBody(t *testing.T, fields map[string]string) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for k, v := range fields {
		fw, err := w.CreateFormField(k)
		require.NoError(t, err)
		_, err = fw.Write([]byte(v))
		require.NoError(t, err)
	}
	require.NoError(t, w.Close())
	return body, w.Boundary()
}

// Dummy values for tests only — not real credentials.
// Must satisfy the strong_password validator (upper, lower, digit, 8+ chars).
const testPassword = "Testpass1dummy"
const wrongPassword = "Wrongpass1dummy"

func newAuthRouter(h *AuthHandler) *gin.Engine {
	r := gin.New()
	v1 := r.Group("/api/v1")
	a := v1.Group("/auth")
	a.POST("/signup", h.Signup)
	a.POST("/login", h.Login)
	a.GET("/me", h.Me)
	return r
}

func TestAuthHandler_Signup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockAuthService)
		h := NewAuthHandler(svc)
		r := newAuthRouter(h)

		expected := &response.SignupResponse{
			UserResponse: response.UserResponse{ID: 1, Phone: "+84901234567"},
			Token:        "jwt-token",
		}
		svc.On("Signup", mock.Anything, request.SignupRequest{
			Phone:    "+84901234567",
			Password: testPassword,
		}).Return(expected, nil)

		body, boundary := multipartBody(t, map[string]string{
			"phone":    "+84901234567",
			"password": testPassword,
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", body)
		req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)
		var resp struct {
			Success bool `json:"success"`
		}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.True(t, resp.Success)
	})

	t.Run("validation fails on missing phone", func(t *testing.T) {
		svc := new(mocks.MockAuthService)
		h := NewAuthHandler(svc)
		r := newAuthRouter(h)

		body, boundary := multipartBody(t, map[string]string{
			"password": testPassword,
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", body)
		req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		svc.AssertNotCalled(t, "Signup")
	})
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockAuthService)
		h := NewAuthHandler(svc)
		r := newAuthRouter(h)

		expected := &response.LoginResponse{
			UserResponse: response.UserResponse{ID: 1, Phone: "+84901234567"},
			Token:        "jwt-token",
		}
		svc.On("Login", mock.Anything, request.LoginRequest{
			Phone:    "+84901234567",
			Password: testPassword,
		}).Return(expected, nil)

		body, boundary := multipartBody(t, map[string]string{
			"phone":    "+84901234567",
			"password": testPassword,
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
		req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid credentials returns 401", func(t *testing.T) {
		svc := new(mocks.MockAuthService)
		h := NewAuthHandler(svc)
		r := newAuthRouter(h)

		svc.On("Login", mock.Anything, mock.Anything).Return(nil, apperror.Unauthorized("auth.invalid_credentials"))

		body, boundary := multipartBody(t, map[string]string{
			"phone":    "+84901234567",
			"password": testPassword,
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
		req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestAuthHandler_Me(t *testing.T) {
	svc := new(mocks.MockAuthService)
	h := NewAuthHandler(svc)
	r := newAuthRouter(h)

	expected := &response.UserResponse{ID: 7, Phone: "+84901234567"}
	svc.On("Me", mock.Anything).Return(expected, nil)

	r.Use(func(c *gin.Context) {
		ctx := auth.WithContext(c.Request.Context(), &auth.Context{UserID: 7, Role: "passenger"})
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}
