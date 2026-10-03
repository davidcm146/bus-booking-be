package repository

import (
	"context"

	"github.com/davidcm146/bus-booking-be/internal/model"
)

// RoleRepository reads the roles reference table.
type RoleRepository interface {
	FindByName(ctx context.Context, name model.UserRole) (*model.Role, error)
}
