package auth

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TokenRepository interface {
	Create(ctx context.Context, rt *RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllByUserID(ctx context.Context, userID string) error
}
type tokenRepository struct {
	db *sqlx.DB
}

// NewTokenRepository recibe *sql.DB y lo envuelve internamente con sqlx.
// Así el caller (wire.go) no se entera de sqlx y es compatible con user.NewUserRepository.
func NewTokenRepository(db *sql.DB) TokenRepository {
	return &tokenRepository{db: sqlx.NewDb(db, "postgres")}
}

// FindByHash — sqlx.GetContext mapea fila a struct automáticamente
func (r *tokenRepository) FindByHash(ctx context.Context, hash string) (*RefreshToken, error) {
	var rt RefreshToken
	err := r.db.GetContext(ctx, &rt,
		`SELECT id, user_id, token_hash, expires_at, created_at, revoked_at
         FROM refresh_tokens WHERE token_hash = $1`, hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

// Create — ExecContext igual que database/sql
func (r *tokenRepository) Create(ctx context.Context, rt *RefreshToken) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
         VALUES ($1, $2, $3, $4)`,
		rt.ID, rt.UserID, rt.TokenHash, rt.ExpiresAt)
	return err
}

// Revoke — UPDATE con named params
func (r *tokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1`, id)
	return err
}

// RevokeAllByUserID — revocación masiva
func (r *tokenRepository) RevokeAllByUserID(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW()
         WHERE user_id = $1 AND revoked_at IS NULL`, userID)
	return err
}
