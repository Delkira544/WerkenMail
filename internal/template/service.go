package template

import (
	"context"

	"github.com/Delkira544/rakiduam/internal/platform/logger"
	"github.com/Delkira544/rakiduam/internal/shared/auth"
	"github.com/Delkira544/rakiduam/internal/shared/emailtemplate"
	"github.com/Delkira544/rakiduam/internal/shared/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	CreateTemplate(ctx context.Context, identity auth.Identity, projectID uuid.UUID, req *CreateTemplateRequest) (*TemplateResponse, error)
}

type ProjectAdapter interface {
	GetProjectByID(ctx context.Context, projectID uuid.UUID) (*ProjectResponse, error)
}

type service struct {
	repo       Repository
	projectSvc ProjectAdapter
}

func NewService(repo Repository, projectSvc ProjectAdapter) Service {
	return &service{
		repo:       repo,
		projectSvc: projectSvc,
	}
}

func (s *service) CreateTemplate(ctx context.Context, identity auth.Identity, projectID uuid.UUID, req *CreateTemplateRequest) (*TemplateResponse, error) {
	log := logger.FromContext(ctx)
	vars, err := toTemplateVariableRequests(req.Variables)
	if err != nil {
		log.Error("Failed to convert template variables", zap.Error(err))
		return nil, errors.BadRequest("Failed to convert template variables: " + err.Error())
	}
	if identity.IsStudent() && identity.UserID != projectID {
		log.Error("Unauthorized access to project", zap.String("user_id", identity.UserID.String()), zap.String("project_id", projectID.String()))
		return nil, errors.Unauthorized("You do not have permission to create a template for this project")
	}

	template := emailtemplate.New(*req.BodyHtml, vars)

	err = template.Validate()
	if err != nil {
		log.Error("Template validation failed", zap.Error(err))
		return nil, errors.BadRequest("Template validation failed: " + err.Error())
	}

	templ := &Template{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      req.Name,
		Subject:   req.Subject,
		BodyHtml:  req.BodyHtml,
		BodyText:  req.BodyText,
		Version:   1,
	}

	varsModel := make([]TemplateVariable, len(vars))

	for i, v := range vars {
		varsModel[i] = TemplateVariable{
			ID:           uuid.New(),
			Key:          v.Key,
			Type:         string(v.Type),
			Required:     v.Required,
			DefaultValue: v.DefaultValue,
		}
	}

	err = s.repo.Create(ctx, templ, varsModel)

	return &TemplateResponse{}, nil
}
