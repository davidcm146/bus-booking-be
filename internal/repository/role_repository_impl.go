package repository

import (
	"context"
	"errors"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"gorm.io/gorm"
)

// roleRepositoryImpl is the GORM-backed implementation of RoleRepository.
type roleRepositoryImpl struct {
	db *gorm.DB
}

// NewRoleRepository creates a new GORM-backed RoleRepository.
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepositoryImpl{db: db}
}

func (r *roleRepositoryImpl) FindByName(ctx context.Context, name model.UserRole) (*model.Role, error) {
	var role model.Role
	if err := r.db.WithContext(ctx).First(&role, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound(i18n.KeyRoleNotFound)
		}
		return nil, apperror.Internal(i18n.KeyRoleFailedFind)
	}
	return &role, nil
}
