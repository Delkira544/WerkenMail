package user

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
	Upsert(ctx context.Context, user *User) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (id, username, name, role, email)
			  VALUES (:id, :username, :name, :role, :email)`
	_, err := r.db.NamedExecContext(ctx, query, user)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	query := `SELECT id, username, name, role,  email FROM users WHERE username = $1`
	err := r.db.GetContext(ctx, &user, query, username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No se encontró el usuario
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Upsert(ctx context.Context, user *User) error {
	query := `INSERT INTO users (id, username, name, role, email)
			  VALUES (:id, :username, :name, :role, :email)
			  ON CONFLICT (username)
			  DO UPDATE SET
				name = EXCLUDED.name,
				role = EXCLUDED.role,
			  	email = EXCLUDED.email`
	_, err := r.db.NamedExecContext(ctx, query, user)

	if err != nil {
		return err
	}

	return nil
}
