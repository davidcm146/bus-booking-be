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
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/davidcm146/bus-booking-be/internal/utils"
)

type AuthServiceImpl struct {
	userRepo  repository.UserRepository
	jwtConfig configs.JWTConfig
}

func NewAuthService(userRepo repository.UserRepository, jwtConfig configs.JWTConfig) *AuthServiceImpl {
	return &AuthServiceImpl{userRepo: userRepo, jwtConfig: jwtConfig}
}

func (s *AuthServiceImpl) Signup(ctx context.Context, req request.SignupRequest) (*response.SignupResponse, error) {
	// check if phone exists
	existing, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err == nil && existing != nil {
		return nil, apperror.Conflict(constant.MsgKeyUserPhoneConflict)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, apperror.Internal(constant.MsgKeyFailedHashPassword)
	}

	user := &model.User{
		Phone:    req.Phone,
		Password: hashedPassword,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, apperror.Internal(constant.MsgKeyFailedRegister)
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role), s.jwtConfig.Secret, s.jwtConfig.ExpiresIn)
	if err != nil {
		return nil, apperror.Internal(constant.MsgKeyFailedGenerateToken)
	}

	return &response.SignupResponse{
		UserResponse: response.ToUserResponse(user),
		Token:        token,
	}, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error) {
	user, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil || user == nil {
		return nil, apperror.Unauthorized(constant.MsgKeyInvalidCredentials)
	}

	if !utils.VerifyPassword(req.Password, user.Password) {
		return nil, apperror.Unauthorized(constant.MsgKeyInvalidCredentials)
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role), s.jwtConfig.Secret, s.jwtConfig.ExpiresIn)
	if err != nil {
		return nil, apperror.Internal(constant.MsgKeyFailedGenerateToken)
	}

	return &response.LoginResponse{
		UserResponse: response.ToUserResponse(user),
		Token:        token,
	}, nil
}

func (s *AuthServiceImpl) Me(ctx context.Context) (*response.UserResponse, error) {
	authCtx, ok := auth.FromContext(ctx)
	if !ok {
		return nil, apperror.Unauthorized(constant.MsgKeyUnauthorized)
	}

	user, err := s.userRepo.FindByID(ctx, authCtx.UserID)
	if err != nil {
		return nil, apperror.Internal(constant.MsgKeyFailedFindUser)
	}

	resp := response.ToUserResponse(user)
	return &resp, nil
}
