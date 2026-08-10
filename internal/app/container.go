package app

import (
	"github.com/davidcm146/bus-booking-be/configs"
	"github.com/davidcm146/bus-booking-be/internal/module"
	"github.com/davidcm146/bus-booking-be/internal/router"
	"gorm.io/gorm"
)

// Container holds all modules and provides handler references for routing.
type Container struct {
	user     *module.UserModule
	auth     *module.AuthModule
	handlers *router.Handlers
}

// NewContainer builds the DI graph: each module wires
func NewContainer(db *gorm.DB, cfg *configs.Config) *Container {
	c := &Container{
		user: module.NewUserModule(db),
		auth: module.NewAuthModule(db, cfg.JWT, cfg.OAuth),
	}

	c.handlers = &router.Handlers{
		User: c.user.Handler,
		Auth: c.auth.Handler,
	}

	return c
}

func (c *Container) Handlers() *router.Handlers {
	return c.handlers
}
