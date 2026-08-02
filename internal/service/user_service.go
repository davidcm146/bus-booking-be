package service

import (
	"context"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/dto/response"
)

// UserService defines the business-logic contract for User operations.
// It accepts request DTOs and returns response DTOs — domain models
// never leak to the handler layer.
type UserService interface {
	CreateUser(ctx context.Context, req request.CreateUserRequest) (*response.UserResponse, error)
	GetUserByID(ctx context.Context, id int) (*response.UserResponse, error)
	GetAllUsers(ctx context.Context) ([]response.UserResponse, error)
	UpdateUser(ctx context.Context, id int, req request.UpdateUserRequest) (*response.UserResponse, error)
	DeleteUser(ctx context.Context, id int) error
}
