package response

import (
	"errors"
	"net/http"

	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/gin-gonic/gin"
)

// APIResponse is the standard envelope for all JSON responses.
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

// Success sends a 200 OK response with data.
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Created sends a 201 Created response with data.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a 204 No Content response.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error sends an error response. If the error is an *apperror.AppError,
// it translates the message key using the request locale;
// otherwise it falls back to a generic localized internal error.
func Error(c *gin.Context, err error) {
	ctx := c.Request.Context()

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPStatus, APIResponse{
			Success: false,
			Error: gin.H{
				"code":    appErr.Code,
				"message": i18n.T(ctx, appErr.MessageKey),
			},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Error: gin.H{
			"code":    constant.CodeInternalError,
			"message": i18n.T(ctx, constant.MsgKeyInternalError),
		},
	})
}

// ValidationError sends a 422 response with field-level validation errors.
func ValidationError(c *gin.Context, details any) {
	ctx := c.Request.Context()

	c.JSON(http.StatusUnprocessableEntity, APIResponse{
		Success: false,
		Error: gin.H{
			"code":    constant.CodeValidationError,
			"message": i18n.T(ctx, constant.MsgKeyValidationFailed),
			"details": details,
		},
	})
}
