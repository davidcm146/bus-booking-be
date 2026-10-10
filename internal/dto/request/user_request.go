package request

// CreateUserRequest is the inbound DTO for creating a new user.
type CreateUserRequest struct {
	FullName string `json:"fullName" binding:"required,min=2,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Role     string `json:"role" binding:"omitempty,oneof=passenger operator admin"`
}

// UpdateUserRequest is the inbound DTO for updating an existing user.
type UpdateUserRequest struct {
	FullName string `json:"fullName" binding:"omitempty,min=2,max=255"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Role     string `json:"role" binding:"omitempty,oneof=passenger operator admin"`
}
