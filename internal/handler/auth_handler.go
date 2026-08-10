package handler

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/davidcm146/bus-booking-be/configs"
	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/service"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/davidcm146/bus-booking-be/internal/shared/response"
	"github.com/davidcm146/bus-booking-be/internal/shared/validator"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// AuthHandler handles HTTP requests for authentication operations.
type AuthHandler struct {
	authService service.AuthService
	oauthConfig configs.OAuthConfig
}

// NewAuthHandler creates a new AuthHandler with the given service.
func NewAuthHandler(authService service.AuthService, oauthConfig configs.OAuthConfig) *AuthHandler {
	return &AuthHandler{authService: authService, oauthConfig: oauthConfig}
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

// GoogleRedirect godoc
// @Summary      Redirect to Google OAuth consent screen
// @Tags         auth
// @Produce      json
// @Success      302
// @Router       /api/v1/auth/google [get]
func (h *AuthHandler) GoogleRedirect(c *gin.Context) {
	oauth2Config := &oauth2.Config{
		ClientID:     h.oauthConfig.GoogleClientID,
		ClientSecret: h.oauthConfig.GoogleClientSecret,
		RedirectURL:  h.oauthConfig.GoogleRedirectURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	url := oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(302, url)
}

// GoogleCallback godoc
// @Summary      Handle Google OAuth callback
// @Tags         auth
// @Produce      json
// @Param        code   query  string  true  "Authorization code"
// @Param        state  query  string  true  "State parameter"
// @Success      200
// @Failure      401
// @Router       /api/v1/auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.Error(c, apperror.BadRequest(constant.MsgKeyBadRequest))
		return
	}

	result, err := h.authService.GoogleOAuth(c.Request.Context(), code)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}
