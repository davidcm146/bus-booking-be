package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/auth"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/davidcm146/bus-booking-be/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const mwSecret = "mw-test-secret"

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	i18n.Init()
	m.Run()
}

// guarded sets up a route protected by AuthMiddleware so we can assert
// whether c.Next() was reached (200) or aborted (4xx).
func guarded(secret string) *gin.Engine {
	r := gin.New()
	r.GET("/protected", AuthMiddleware(secret), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func TestAuthMiddleware(t *testing.T) {
	token, err := utils.GenerateToken(5, "passenger", mwSecret, time.Hour)
	require.NoError(t, err)

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{"valid token", "Bearer " + token, http.StatusOK},
		{"missing header", "", http.StatusUnauthorized},
		{"wrong scheme", "Basic " + token, http.StatusUnauthorized},
		{"invalid token", "Bearer garbage-token", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := guarded(mwSecret)

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	r := guarded(mwSecret)
	token, err := utils.GenerateToken(5, "passenger", "different-secret", time.Hour)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_SetsAuthContext(t *testing.T) {
	r := gin.New()
	r.GET("/check", AuthMiddleware(mwSecret), func(c *gin.Context) {
		authCtx, ok := auth.FromContext(c.Request.Context())
		require.True(t, ok)
		require.Equal(t, 5, authCtx.UserID)
		require.Equal(t, "passenger", authCtx.Role)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	token, err := utils.GenerateToken(5, "passenger", mwSecret, time.Hour)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/check", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireRole(t *testing.T) {
	adminToken, err := utils.GenerateToken(1, "admin", mwSecret, time.Hour)
	require.NoError(t, err)
	passengerToken, err := utils.GenerateToken(2, "passenger", mwSecret, time.Hour)
	require.NoError(t, err)

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{"allowed role", "Bearer " + adminToken, http.StatusOK},
		{"other role is forbidden", "Bearer " + passengerToken, http.StatusForbidden},
		{"no token is unauthorized", "", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/admin", AuthMiddleware(mwSecret), RequireRole(model.RoleAdmin), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"ok": true})
			})

			req := httptest.NewRequest(http.MethodGet, "/admin", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestRequireRole_WithoutAuthMiddleware(t *testing.T) {
	r := gin.New()
	r.GET("/admin", RequireRole(model.RoleAdmin), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin", nil))

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
