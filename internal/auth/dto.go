package auth

import "github.com/google/uuid"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_at"`
	Role         string `json:"role"`
}

// DTO for adapters
type SyncUserRequest struct {
	Username string
	Name     string
	Email    string
	Role     string
}

type SyncUserResponse struct {
	ID       uuid.UUID
	Username string
	Role     string
}
