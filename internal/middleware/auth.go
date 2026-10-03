package middleware

import (
	"slices"
	"strings"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/auth"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/davidcm146/bus-booking-be/internal/shared/response"
	"github.com/davidcm146/bus-booking-be/internal/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token from the Authorization header
// and sets userID and role in the gin context for downstream handlers.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, apperror.Unauthorized(i18n.KeyAuthMissingHeader))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, apperror.Unauthorized(i18n.KeyAuthInvalidFormat))
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1], jwtSecret)
		if err != nil {
			response.Error(c, apperror.Unauthorized(i18n.KeyAuthInvalidToken))
			c.Abort()
			return
		}

		authCtx := &auth.Context{
			UserID: claims.UserID,
			Role:   claims.Role,
		}

		ctx := auth.WithContext(c.Request.Context(), authCtx)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// RequireRole allows the request through only when the authenticated user
// has one of the given roles. It must run after AuthMiddleware.
func RequireRole(roles ...model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCtx, ok := auth.FromContext(c.Request.Context())
		if !ok {
			response.Error(c, apperror.ErrUnauthorized)
			c.Abort()
			return
		}
		if !slices.Contains(roles, model.UserRole(authCtx.Role)) {
			response.Error(c, apperror.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
