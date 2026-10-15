package param

type LoginRequest struct {
	Email string `json:"email"`
}

// TODO: add user usage field
type LoginResponse struct {
	ApiKey string `json:"api_key"`
	Token  int64  `json:"token"`
}
