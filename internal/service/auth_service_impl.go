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
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/davidcm146/bus-booking-be/internal/utils"
	"golang.org/x/oauth2"
)

type AuthServiceImpl struct {
	userRepo    repository.UserRepository
	roleRepo    repository.RoleRepository
	jwtConfig   configs.JWTConfig
	oauthConfig configs.OAuthConfig
}

func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, jwtConfig configs.JWTConfig, oauthConfig configs.OAuthConfig) *AuthServiceImpl {
	return &AuthServiceImpl{userRepo: userRepo, roleRepo: roleRepo, jwtConfig: jwtConfig, oauthConfig: oauthConfig}
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

func (s *AuthServiceImpl) GoogleOAuth(ctx context.Context, code string) (*response.LoginResponse, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, apperror.Internal(i18n.KeyAuthFailedHashPassword)
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
		return nil, apperror.Unauthorized(i18n.KeyAuthFailedOAuth)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, apperror.Unauthorized(i18n.KeyAuthFailedOAuth)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: s.oauthConfig.GoogleClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, apperror.Unauthorized(i18n.KeyAuthFailedOAuth)
	}

	var claims request.GoogleClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, apperror.Internal(i18n.KeyErrorInternal)
	}

	// Find or create user by email
	user, err := s.userRepo.FindByEmail(ctx, claims.Email)
	if err != nil || user == nil {
		role, err := s.roleRepo.FindByName(ctx, model.RolePassenger)
		if err != nil {
			return nil, err
		}

		user = &model.User{
			FullName: claims.Name,
			Email:    &claims.Email,
			Role:     role,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, apperror.Internal(i18n.KeyAuthFailedRegister)
		}
	}

	jwtToken, err := utils.GenerateToken(user.ID, string(user.RoleName()), s.jwtConfig.Secret, s.jwtConfig.ExpiresIn)
	if err != nil {
		return nil, apperror.Internal(i18n.KeyAuthFailedGenerateToken)
	}

	return &response.LoginResponse{
		UserResponse: response.ToUserResponse(user),
		Token:        jwtToken,
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
