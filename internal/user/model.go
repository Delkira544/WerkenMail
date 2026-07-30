package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `db:"id"`
	Username  string     `db:"username"`
	Name      string     `db:"name"`
	Role      string     `db:"role"`
	Email     string     `db:"email"`
	CreatedAt *time.Time `db:"created_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID.String(),
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.String(),
	}
}
