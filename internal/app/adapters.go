package app

import (
	"context"

	"github.com/Delkira544/rakiduam/internal/auth"
	"github.com/Delkira544/rakiduam/internal/user"
)

type authUserAdapter struct {
	svc user.UserService
}

func (a *authUserAdapter) SyncFromLDAP(ctx context.Context, req *auth.SyncUserRequest) error {
	return a.svc.SyncFromLDAP(ctx, &user.CreateUserRequest{
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
	})
}
