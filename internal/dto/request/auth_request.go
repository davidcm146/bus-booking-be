package request

type SignupRequest struct {
	Phone    string `form:"phone" binding:"required,phone"`
	Password string `form:"password" binding:"required,strong_password"`
}

type LoginRequest struct {
	Phone    string `form:"phone" binding:"required,phone"`
	Password string `form:"password" binding:"required"`
}

type GoogleClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
