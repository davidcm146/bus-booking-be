package service

import (
	"context"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/dto/response"
	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/repository"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/davidcm146/bus-booking-be/internal/utils"
)

// userServiceImpl implements UserService.
type userServiceImpl struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

// NewUserService creates a new UserService with the given repositories.
func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository) UserService {
	return &userServiceImpl{userRepo: userRepo, roleRepo: roleRepo}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req request.CreateUserRequest) (*response.UserResponse, error) {
	roleName := model.RolePassenger
	if req.Role != "" {
		roleName = model.UserRole(req.Role)
	}
	role, err := s.roleRepo.FindByName(ctx, roleName)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, apperror.Internal(i18n.KeyAuthFailedHashPassword)
	}

	user := &model.User{
		FullName: req.FullName,
		Email:    &req.Email,
		Phone:    &req.Phone,
		Password: hashedPassword,
		RoleID:   &role.ID,
		Role:     role,
		Status:   model.UserStatusActive,
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
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	if req.Role != "" {
		role, err := s.roleRepo.FindByName(ctx, model.UserRole(req.Role))
		if err != nil {
			return nil, err
		}
		user.RoleID, user.Role = &role.ID, role
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
