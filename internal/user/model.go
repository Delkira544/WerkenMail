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
