package user

import (
	"time"

	authshared "github.com/Delkira544/rakiduam/internal/shared/auth"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID       `db:"id"`
	Username  string          `db:"username"`
	Name      string          `db:"name"`
	Role      authshared.Role `db:"role"`
	Email     string          `db:"email"`
	CreatedAt *time.Time      `db:"created_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID.String(),
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role.String(),
		CreatedAt: u.CreatedAt.String(),
	}
}
