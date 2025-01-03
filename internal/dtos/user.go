package dtos

import "time"

type CreateUserRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CreateUserResponse struct {
	UserID      uint      `json:"user_id"`
	AccessToken string    `json:"access_token"`
	ExpiredAt   time.Time `json:"expired_at"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `binding:"required" json:"user_name"`
	Password string `binding:"required" json:"password"`
}

type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiredAt   time.Time `json:"expired_at"`
}
