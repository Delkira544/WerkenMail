package app

import (
	"context"

	"github.com/Delkira544/rakiduam/internal/auth"
	"github.com/Delkira544/rakiduam/internal/projects"
	"github.com/Delkira544/rakiduam/internal/template"
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

type projectTemplateAdapter struct {
	svc projects.Service
}

func (a *projectTemplateAdapter) GetProjectByID(ctx context.Context, projectID uuid.UUID) (*template.ProjectResponse, error) {
	p, err := a.svc.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &template.ProjectResponse{
		ID:          p.ID,
		UserID:      p.UserID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}
