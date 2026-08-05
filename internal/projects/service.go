package projects

import (
	"context"

	"github.com/Delkira544/rakiduam/internal/shared/auth"
	"github.com/Delkira544/rakiduam/internal/shared/errors"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, input *CreateProjectInput) (*ProjectResponse, error)
	List(ctx context.Context, identity auth.Identity, req ListProjectRequest) ([]ProjectResponse, error)
	Update(ctx context.Context, id uuid.UUID, identity auth.Identity, req UpdateProjectRequest) (*ProjectResponse, error)
	Delete(ctx context.Context, id uuid.UUID, identity auth.Identity) error
	GetProjectByID(ctx context.Context, id uuid.UUID) (*ProjectResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, input *CreateProjectInput) (*ProjectResponse, error) {
	var project = &Project{
		ID:          uuid.New(),
		UserID:      input.UserID,
		Name:        input.Name,
		Description: input.Description,
	}
	err := s.repo.Create(ctx, project)
	if err != nil {
		return nil, errors.Internal("create a project")
	}
	return project.ToResponse(), nil
}

func (s *service) GetProjectByID(ctx context.Context, id uuid.UUID) (*ProjectResponse, error) {
	project, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, errors.Internal("get project by id")
	}
	if project == nil {
		return nil, errors.NotFound("project not found")
	}
	return project.ToResponse(), nil
}

func (s *service) List(ctx context.Context, identity auth.Identity, req ListProjectRequest) ([]ProjectResponse, error) {
	filter := ProjectFilter{}

	if identity.IsStudent() {
		uid := identity.UserID.String()
		filter.UserID = &uid
	}

	if req.Search != nil {
		filter.Search = req.Search
	}
	if req.Limit != nil {
		filter.Limit = req.Limit
	}
	if req.Offset != nil {
		filter.Offset = req.Offset
	}
	projects, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, errors.Internal("list projects")
	}
	return toResponse(projects), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, identity auth.Identity, req UpdateProjectRequest) (*ProjectResponse, error) {
	project, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, errors.Internal("get project by id")
	}
	if project == nil {
		return nil, errors.NotFound("project not found")
	}
	if identity.IsStudent() && project.UserID.String() != identity.UserID.String() {
		return nil, errors.Forbidden("not the owner")
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = *req.Description
	}

	if err := s.repo.Update(ctx, project); err != nil {
		return nil, errors.Internal("update project")
	}
	return project.ToResponse(), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, identity auth.Identity) error {
	project, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return errors.Internal("get project by id")
	}
	if project == nil {
		return errors.NotFound("project not found")
	}
	if identity.IsStudent() && project.UserID.String() != identity.UserID.String() {
		return errors.Forbidden("not the owner")
	}

	if err := s.repo.Delete(ctx, id.String()); err != nil {
		return errors.Internal("delete project")
	}
	return nil
}
