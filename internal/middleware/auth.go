package middleware

import (
	"strings"

	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/auth"
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
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
			response.Error(c, apperror.Unauthorized(constant.MsgKeyMissingAuthHeader))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, apperror.Unauthorized(constant.MsgKeyInvalidAuthFormat))
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1], jwtSecret)
		if err != nil {
			response.Error(c, apperror.Unauthorized(constant.MsgKeyInvalidToken))
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
