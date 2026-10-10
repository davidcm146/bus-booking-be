package response

import (
	"time"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/utils"
)

// UserResponse is the outbound DTO for user data.
// It never exposes sensitive fields like password.
type UserResponse struct {
	ID        int            `json:"id"`
	FullName  string         `json:"fullName"`
	Email     string         `json:"email"`
	Phone     string         `json:"phone"`
	Role      model.UserRole `json:"role"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// ToUserResponse maps a domain model to a response DTO.
func ToUserResponse(u *model.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		FullName:  u.FullName,
		Email:     utils.DeferString(u.Email),
		Phone:     utils.DeferString(u.Phone),
		Role:      u.RoleName(),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToUserResponses maps a slice of domain models to response DTOs.
func ToUserResponses(users []model.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, u := range users {
		responses[i] = ToUserResponse(&u)
	}
	return responses
}
