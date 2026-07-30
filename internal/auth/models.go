package auth

import (
	"time"

	authshared "github.com/Delkira544/rakiduam/internal/shared/auth"
	"github.com/google/uuid"
)

func ToRole(gidNumber string) authshared.Role {
	switch gidNumber {
	case "600":
		return authshared.RoleStudent
	case "500":
		return authshared.RoleFunc
	default:
		return authshared.RoleStudent // Default role if not recognized
	}
}

type User struct {
	Username string
	FullName string
	Email    string
	DN       string
	Role     authshared.Role
}

type RefreshToken struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	TokenHash string     `db:"token_hash"`
	ExpiresAt time.Time  `db:"expires_at"`
	CreatedAt time.Time  `db:"created_at"`
	RevokedAt *time.Time `db:"revoked_at"`
}
