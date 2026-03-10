package models

// LoginRequest - Note: Capital letter = exported/"public"
// Small letters = private to package

type LoginRequest struct {
	Username string `json:"username"` // json tag = JSON serialisation hint
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"` // omitempty = exclude if empty
}

type Claims struct {
	Username  string `json:"username"`
	ExpiresAt int64  `json:"exp"`
}
