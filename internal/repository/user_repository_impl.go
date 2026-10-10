package repository

import (
	"context"
	"errors"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/dberror"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// userConflicts maps unique constraints on users to their conflict messages.
var userConflicts = []dberror.Option{
	dberror.WithConflict("users_email_key", i18n.KeyUserEmailConflict),
	dberror.WithConflict("users_phone_key", i18n.KeyUserPhoneConflict),
}

type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM-backed UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) active(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Joins("Role").Where("users.deleted = ?", false)
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Omit(clause.Associations).Create(user).Error; err != nil {
		return dberror.Translate(err, i18n.KeyUserFailedCreate, userConflicts...)
	}
	return nil
}

func (r *userRepositoryImpl) findOne(ctx context.Context, query string, arg any) (*model.User, error) {
	var user model.User
	if err := r.active(ctx).First(&user, query, arg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound(i18n.KeyUserNotFound)
		}
		return nil, dberror.Translate(err, i18n.KeyUserFailedFind)
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id int) (*model.User, error) {
	return r.findOne(ctx, "users.id = ?", id)
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.findOne(ctx, "users.email = ?", email)
}

func (r *userRepositoryImpl) FindByPhone(ctx context.Context, phoneNumber string) (*model.User, error) {
	return r.findOne(ctx, "users.phone = ?", phoneNumber)
}

func (r *userRepositoryImpl) FindAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := r.active(ctx).Find(&users).Error; err != nil {
		return nil, dberror.Translate(err, i18n.KeyUserFailedFetch)
	}
	return users, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Model(user).Where("deleted = ?", false).
		Select("full_name", "email", "phone", "role_id", "operator_id", "language", "status", "updated_at").
		Updates(user)
	if result.Error != nil {
		return dberror.Translate(result.Error, i18n.KeyUserFailedUpdate, userConflicts...)
	}
	if result.RowsAffected == 0 {
		return apperror.NotFound(i18n.KeyUserNotFound)
	}
	return nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ? AND deleted = ?", id, false).
		Updates(map[string]any{"deleted": true, "deleted_at": gorm.Expr("CURRENT_TIMESTAMP")})
	if result.Error != nil {
		return dberror.Translate(result.Error, i18n.KeyUserFailedDelete)
	}
	if result.RowsAffected == 0 {
		return apperror.NotFound(i18n.KeyUserNotFound)
	}
	return nil
}
