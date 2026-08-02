package router

import (
	"github.com/davidcm146/bus-booking-be/internal/handler"
	"github.com/davidcm146/bus-booking-be/internal/middleware"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/gin-gonic/gin"
)

// Handlers groups all handler references the router needs.
type Handlers struct {
	User *handler.UserHandler
	Auth *handler.AuthHandler
}

// New creates and configures the Gin engine with all route groups.
func New(h *Handlers, i18nService i18n.I18n, jwtSecret string) *gin.Engine {
	r := gin.Default()

	// --- Global middleware ---
	r.Use(middleware.LocaleMiddleware(i18nService))

	// --- Health check ---
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// --- API v1 ---
	v1 := r.Group("/api/v1")
	{
		registerUserRoutes(v1, h.User)
		registerAuthRoutes(v1, h.Auth, jwtSecret)
	}

	return r
}

func registerAuthRoutes(rg *gin.RouterGroup, h *handler.AuthHandler, jwtSecret string) {
	auth := rg.Group("/auth")
	{
		auth.POST("/signup", h.Signup)
		auth.POST("/login", h.Login)
		auth.GET("/me", middleware.AuthMiddleware(jwtSecret), h.Me)
	}
}

// registerUserRoutes registers all /users endpoints.
func registerUserRoutes(rg *gin.RouterGroup, h *handler.UserHandler) {
	users := rg.Group("/users")
	{
		users.POST("", h.Create)
		users.GET("", h.GetAll)
		users.GET("/:id", h.GetByID)
		users.PUT("/:id", h.Update)
		users.DELETE("/:id", h.Delete)
	}
}
