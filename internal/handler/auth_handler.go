package handler

import (
	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/service"
	"github.com/davidcm146/bus-booking-be/internal/shared/response"
	"github.com/davidcm146/bus-booking-be/internal/shared/validator"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication operations.
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new AuthHandler with the given service.
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Signup godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       mpfd
// @Produce      json
// @Param        name      formData  string  true   "User name"
// @Param        email     formData  string  true   "Email address"
// @Param        phone     formData  string  true   "Phone number"
// @Param        password  formData  string  true   "Password"
// @Param        role      formData  string  false  "User role"
// @Success      201
// @Failure      400
// @Failure      409
// @Failure      422
// @Router       /api/v1/auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req request.SignupRequest
	if err := c.ShouldBind(&req); err != nil {
		response.ValidationError(c, validator.FormatError(c.Request.Context(), err))
		return
	}

	result, err := h.authService.Signup(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, result)
}

// Login godoc
// @Summary      Authenticate a user
// @Tags         auth
// @Accept       mpfd
// @Produce      json
// @Param        phone     formData  string  true  "Phone number"
// @Param        password  formData  string  true  "Password"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      422
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		response.ValidationError(c, validator.FormatError(c.Request.Context(), err))
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// Me godoc
// @Summary      Get current authenticated user
// @Tags         auth
// @Accept       mpfd
// @Produce      json
// @Success      200
// @Failure      401
// @Failure      422
// @Router       /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	result, err := h.authService.Me(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}
