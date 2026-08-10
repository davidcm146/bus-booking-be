package service

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/davidcm146/bus-booking-be/configs"
	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/dto/response"
	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/repository"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/auth"
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/davidcm146/bus-booking-be/internal/utils"
	"golang.org/x/oauth2"
)

type AuthServiceImpl struct {
	userRepo    repository.UserRepository
	jwtConfig   configs.JWTConfig
	oauthConfig configs.OAuthConfig
}

func NewAuthService(userRepo repository.UserRepository, jwtConfig configs.JWTConfig, oauthConfig configs.OAuthConfig) *AuthServiceImpl {
	return &AuthServiceImpl{userRepo: userRepo, jwtConfig: jwtConfig, oauthConfig: oauthConfig}
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

func (s *AuthServiceImpl) GoogleOAuth(ctx context.Context, code string) (*response.LoginResponse, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, apperror.Internal(constant.MsgKeyInternalError)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     s.oauthConfig.GoogleClientID,
		ClientSecret: s.oauthConfig.GoogleClientSecret,
		RedirectURL:  s.oauthConfig.GoogleRedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, apperror.Unauthorized(constant.MsgKeyInvalidCredentials)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, apperror.Unauthorized(constant.MsgKeyInvalidCredentials)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: s.oauthConfig.GoogleClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, apperror.Unauthorized(constant.MsgKeyInvalidCredentials)
	}

	var claims request.GoogleClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, apperror.Internal(constant.MsgKeyInternalError)
	}

	// Find or create user by email
	user, err := s.userRepo.FindByEmail(ctx, claims.Email)
	if err != nil || user == nil {
		user = &model.User{
			FullName: claims.Name,
			Email:    claims.Email,
			Role:     model.RolePassenger,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, apperror.Internal(constant.MsgKeyFailedRegister)
		}
	}

	jwtToken, err := utils.GenerateToken(user.ID, string(user.Role), s.jwtConfig.Secret, s.jwtConfig.ExpiresIn)
	if err != nil {
		return nil, apperror.Internal(constant.MsgKeyFailedGenerateToken)
	}

	return &response.LoginResponse{
		UserResponse: response.ToUserResponse(user),
		Token:        jwtToken,
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
