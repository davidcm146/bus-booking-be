package repository

import (
	"context"
	"errors"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const pgUniqueViolation = "23505"

// userRepositoryImpl is the GORM-backed implementation of UserRepository.
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM-backed UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// active scopes queries to non-deleted users and loads their role.
func (r *userRepositoryImpl) active(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Joins("Role").Where("users.deleted = ?", false)
}

// writeError maps unique violations to conflicts and anything else to fallback.
func writeError(err error, fallback constant.MsgKey) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != pgUniqueViolation {
		return apperror.Internal(fallback)
	}
	switch pgErr.ConstraintName {
	case "users_email_key":
		return apperror.Conflict(i18n.KeyUserEmailConflict)
	case "users_phone_key":
		return apperror.Conflict(i18n.KeyUserPhoneConflict)
	}
	return apperror.Conflict(i18n.KeyErrorResourceConflict)
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Omit(clause.Associations).Create(user).Error; err != nil {
		return writeError(err, i18n.KeyUserFailedCreate)
	}
	return nil
}

func (r *userRepositoryImpl) findOne(ctx context.Context, query string, arg any) (*model.User, error) {
	var user model.User
	if err := r.active(ctx).First(&user, query, arg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound(i18n.KeyUserNotFound)
		}
		return nil, apperror.Internal(i18n.KeyUserFailedFind)
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
		return nil, apperror.Internal(i18n.KeyUserFailedFetch)
	}
	return users, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Model(user).Where("deleted = ?", false).
		Select("full_name", "email", "phone", "role_id", "operator_id", "language", "status", "updated_at").
		Updates(user)
	if result.Error != nil {
		return writeError(result.Error, i18n.KeyUserFailedUpdate)
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
		return apperror.Internal(i18n.KeyUserFailedDelete)
	}
	if result.RowsAffected == 0 {
		return apperror.NotFound(i18n.KeyUserNotFound)
	}
	return nil
}
