package service

import (
	"context"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/dto/response"
	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/repository"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
)

// userServiceImpl implements UserService.
type userServiceImpl struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService with the given repository.
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req request.CreateUserRequest) (*response.UserResponse, error) {
	// Check for duplicate email
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, apperror.Conflict(constant.MsgKeyUserEmailConflict)
	}

	role := model.RolePassenger
	if req.Role != "" {
		role = model.UserRole(req.Role)
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password, // TODO: hash password before storing
		Role:     role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	resp := response.ToUserResponse(user)
	return &resp, nil
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, id int) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := response.ToUserResponse(user)
	return &resp, nil
}

func (s *userServiceImpl) GetAllUsers(ctx context.Context) ([]response.UserResponse, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return response.ToUserResponses(users), nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, id int, req request.UpdateUserRequest) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply partial updates
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if new email conflicts with another user
		existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
		if existing != nil && existing.ID != id {
			return nil, apperror.Conflict(constant.MsgKeyUserEmailConflict)
		}
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Role != "" {
		user.Role = model.UserRole(req.Role)
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	resp := response.ToUserResponse(user)
	return &resp, nil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, id int) error {
	return s.userRepo.Delete(ctx, id)
}
