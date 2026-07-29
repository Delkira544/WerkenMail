package auth

import (
	"time"

	"github.com/google/uuid"
)

type UserRole int

const (
	RoleStudent UserRole = iota
	RoleFunc
)

func (r UserRole) String() string {
	switch r {
	case RoleStudent:
		return "student"
	case RoleFunc:
		return "func"
	default:
		return "student"
	}
}

func ToRole(roleStr string) UserRole {
	switch roleStr {
	case "600":
		return RoleStudent
	case "500":
		return RoleFunc
	default:
		return RoleStudent // Default role if not recognized
	}
}

type User struct {
	Username string ``
	FullName string
	Email    string
	DN       string
	Role     UserRole
}

type RefreshToken struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	TokenHash string     `db:"token_hash"`
	ExpiresAt time.Time  `db:"expires_at"`
	CreatedAt time.Time  `db:"created_at"`
	RevokedAt *time.Time `db:"revoked_at"`
}
