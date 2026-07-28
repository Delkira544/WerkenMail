package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `db:"id"`
	LDAPUID   string     `db:"ldap_uid"`
	FullName  string     `db:"full_name"`
	Email     string     `db:"email"`
	CreatedAt *time.Time `db:"created_at"`
}
