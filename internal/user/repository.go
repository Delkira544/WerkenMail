package user

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByLDAPUID(ctx context.Context, ldapUID string) (*User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (id, ldap_uid, full_name, email)
			  VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.LDAPUID, user.FullName, user.Email)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) FindByLDAPUID(ctx context.Context, ldapUID string) (*User, error) {
	query := `SELECT id, ldap_uid, full_name, email FROM users WHERE ldap_uid = $1`
	row := r.db.QueryRowContext(ctx, query, ldapUID)

	var user User
	err := row.Scan(&user.ID, &user.LDAPUID, &user.FullName, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No se encontró el usuario
		}
		return nil, err
	}

	return &user, nil

}
