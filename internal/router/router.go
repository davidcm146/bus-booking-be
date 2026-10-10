package router

import (
	"github.com/davidcm146/bus-booking-be/internal/handler"
	"github.com/davidcm146/bus-booking-be/internal/middleware"
	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/gin-gonic/gin"
)

// Handlers groups all handler references the router needs.
type Handlers struct {
	User *handler.UserHandler
	Auth *handler.AuthHandler
	Trip *handler.TripHandler
}

// New creates and configures the Gin engine with all route groups.
func New(h *Handlers, jwtSecret string) *gin.Engine {
	r := gin.Default()

	// --- Global middleware ---
	r.Use(middleware.I18n())

	// --- Health check ---
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// --- API v1 ---

	v1 := r.Group("/api/v1")
	{
		registerUserRoutes(v1, h.User, jwtSecret)
		registerAuthRoutes(v1, h.Auth, jwtSecret)
		registerTripRoutes(v1, h.Trip)
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

// registerUserRoutes registers all /users endpoints (admin only).
func registerUserRoutes(rg *gin.RouterGroup, h *handler.UserHandler, jwtSecret string) {
	users := rg.Group("/users", middleware.AuthMiddleware(jwtSecret), middleware.RequireRole(model.RoleAdmin))
	{
		users.POST("", h.Create)
		users.GET("", h.GetAll)
		users.GET("/:id", h.GetByID)
		users.PUT("/:id", h.Update)
		users.DELETE("/:id", h.Delete)
	}
}

// registerTripRoutes registers all /trips endpoints (public).
func registerTripRoutes(rg *gin.RouterGroup, h *handler.TripHandler) {
	trips := rg.Group("/trips")
	{
		trips.GET("", h.Search)
		trips.GET("/:id", h.GetByID)
	}
}
