package service

import (
	"context"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/dto/response"
)

type AuthService interface {
	Signup(ctx context.Context, req request.SignupRequest) (*response.SignupResponse, error)
	Login(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error)
	Me(ctx context.Context) (*response.UserResponse, error)
}
