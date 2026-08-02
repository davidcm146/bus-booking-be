package repository

import (
	"context"
	"errors"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"gorm.io/gorm"
)

// userRepositoryImpl is the GORM-backed implementation of UserRepository.
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM-backed UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return apperror.Internal(constant.MsgKeyFailedCreateUser)
	}
	return nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound(constant.MsgKeyUserNotFound)
		}
		return nil, apperror.Internal(constant.MsgKeyFailedFindUser)
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound(constant.MsgKeyUserNotFound)
		}
		return nil, apperror.Internal(constant.MsgKeyFailedFindUser)
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByPhone(ctx context.Context, phoneNumber string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "phone = ?", phoneNumber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound(constant.MsgKeyUserNotFound)
		}
		return nil, apperror.Internal(constant.MsgKeyFailedFindUser)
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, apperror.Internal(constant.MsgKeyFailedFetchUsers)
	}
	return users, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return apperror.Internal(constant.MsgKeyFailedUpdateUser)
	}
	return nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return apperror.Internal(constant.MsgKeyFailedDeleteUser)
	}
	if result.RowsAffected == 0 {
		return apperror.NotFound(constant.MsgKeyUserNotFound)
	}
	return nil
}
