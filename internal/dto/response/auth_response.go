package response

type SignupResponse struct {
	UserResponse
	Token string `json:"token"`
}

type LoginResponse struct {
	UserResponse
	Token string `json:"token"`
}
