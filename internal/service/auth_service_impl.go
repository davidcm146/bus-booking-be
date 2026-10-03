package service

import (
	"context"

	"github.com/davidcm146/bus-booking-be/configs"
	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/dto/response"
	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/repository"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/auth"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/davidcm146/bus-booking-be/internal/utils"
)

type AuthServiceImpl struct {
	userRepo  repository.UserRepository
	roleRepo  repository.RoleRepository
	jwtConfig configs.JWTConfig
}

func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, jwtConfig configs.JWTConfig) *AuthServiceImpl {
	return &AuthServiceImpl{userRepo: userRepo, roleRepo: roleRepo, jwtConfig: jwtConfig}
}

func (s *AuthServiceImpl) Signup(ctx context.Context, req request.SignupRequest) (*response.SignupResponse, error) {
	role, err := s.roleRepo.FindByName(ctx, model.RolePassenger)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, apperror.Internal(i18n.KeyAuthFailedHashPassword)
	}

	user := &model.User{
		Phone:    &req.Phone,
		Password: hashedPassword,
		RoleID:   &role.ID,
		Role:     role,
		Status:   model.UserStatusActive,
	}

	// Phone uniqueness is enforced by the DB; the repository maps violations to 409.
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, string(user.RoleName()), s.jwtConfig.Secret, s.jwtConfig.ExpiresIn)
	if err != nil {
		return nil, apperror.Internal(i18n.KeyAuthFailedGenerateToken)
	}

	return &response.SignupResponse{
		UserResponse: response.ToUserResponse(user),
		Token:        token,
	}, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error) {
	user, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil || user == nil {
		return nil, apperror.Unauthorized(i18n.KeyAuthInvalidCredentials)
	}

	if !utils.VerifyPassword(req.Password, user.Password) {
		return nil, apperror.Unauthorized(i18n.KeyAuthInvalidCredentials)
	}

	token, err := utils.GenerateToken(user.ID, string(user.RoleName()), s.jwtConfig.Secret, s.jwtConfig.ExpiresIn)
	if err != nil {
		return nil, apperror.Internal(i18n.KeyAuthFailedGenerateToken)
	}

	return &response.LoginResponse{
		UserResponse: response.ToUserResponse(user),
		Token:        token,
	}, nil
}

func (s *AuthServiceImpl) Me(ctx context.Context) (*response.UserResponse, error) {
	authCtx, ok := auth.FromContext(ctx)
	if !ok {
		return nil, apperror.Unauthorized(i18n.KeyErrorUnauthorized)
	}

	user, err := s.userRepo.FindByID(ctx, authCtx.UserID)
	if err != nil {
		return nil, err
	}

	resp := response.ToUserResponse(user)
	return &resp, nil
}
