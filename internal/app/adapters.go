package app

import (
	"context"

	"github.com/Delkira544/rakiduam/internal/auth"
	"github.com/Delkira544/rakiduam/internal/user"
	"github.com/google/uuid"
)

type authUserAdapter struct {
	svc user.UserService
}

func (a *authUserAdapter) SyncFromLDAP(ctx context.Context, req *auth.SyncUserRequest) (*auth.SyncUserResponse, error) {
	u, err := a.svc.SyncFromLDAP(ctx, &user.CreateUserRequest{
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
	})
	if err != nil {
		return nil, err
	}
	id, _ := uuid.Parse(u.ID)
	return &auth.SyncUserResponse{
		ID:       id,
		Username: u.Username,
		Role:     u.Role,
	}, nil
}
