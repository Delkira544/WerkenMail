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

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTemplate(ctx context.Context, identity auth.Identity, projectID uuid.UUID, req *CreateTemplateRequest) (*TemplateResponse, error) {
	log := logger.FromContext(ctx)
	vars, err := toTemplateVariableRequests(req.Variables)
	if err != nil {
		log.Error("Failed to convert template variables", zap.Error(err))
		return nil, errors.BadRequest("Failed to convert template variables: " + err.Error())
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
