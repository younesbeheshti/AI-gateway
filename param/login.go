package param

type LoginRequest struct {
	Email string `json:"email"`
}

// TODO: add user usage field
type LoginResponse struct {
	Token string `json:"your_token"`
}
