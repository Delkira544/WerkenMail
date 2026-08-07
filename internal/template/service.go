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
	Create(ctx context.Context, identity auth.Identity, projectID uuid.UUID, req *CreateTemplateRequest) (*TemplateResponse, error)
	GetByID(ctx context.Context, identity auth.Identity, templateID uuid.UUID) (*TemplateResponse, error)
	ListByProjectID(ctx context.Context, identity auth.Identity, projectID uuid.UUID) ([]TemplateListItemResponse, error)
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

func (s *service) Create(ctx context.Context, identity auth.Identity, projectID uuid.UUID, req *CreateTemplateRequest) (*TemplateResponse, error) {
	log := logger.FromContext(ctx)
	if _, err := s.verifyProjectAccess(ctx, identity, projectID); err != nil {
		return nil, err
	}

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
			DefaultValue: v.DefaultValue,
		}
	}

	err = s.repo.Create(ctx, templ, varsModel)
	if err != nil {
		log.Error("Failed to create template", zap.Error(err))
		return nil, errors.Internal("Failed to create template: " + err.Error())
	}
	res := make([]TemplateVariableResponse, len(varsModel))
	for i, v := range varsModel {
		res[i] = v.ToResponse()
	}

	return &TemplateResponse{
		ID:                       templ.ID,
		Name:                     templ.Name,
		Subject:                  templ.Subject,
		BodyHtml:                 templ.BodyHtml,
		BodyText:                 templ.BodyText,
		Version:                  templ.Version,
		TemplateVariableResponse: res,
		CreatedAt:                templ.CreatedAt,
		UpdatedAt:                templ.UpdatedAt,
	}, nil
}

func (s *service) ListByProjectID(ctx context.Context, identity auth.Identity, projectID uuid.UUID) ([]TemplateListItemResponse, error) {
	if _, err := s.verifyProjectAccess(ctx, identity, projectID); err != nil {
		return nil, err
	}

	templates, err := s.repo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, errors.Internal("list templates")
	}

	ids := make([]uuid.UUID, len(templates))
	for i, t := range templates {
		ids[i] = t.ID
	}

	counts, err := s.repo.GetVariableCounts(ctx, ids)
	if err != nil {
		return nil, errors.Internal("count variables")
	}

	responses := make([]TemplateListItemResponse, len(templates))
	for i, t := range templates {
		responses[i] = t.ToListItemResponse(counts[t.ID])
	}
	return responses, nil
}

func (s *service) GetByID(ctx context.Context, identity auth.Identity, templateID uuid.UUID) (*TemplateResponse, error) {
	template, err := s.repo.GetByID(ctx, templateID)
	if err != nil {
		return nil, errors.Internal("get template by id")
	}
	if template == nil {
		return nil, errors.NotFound("template not found")
	}

	if _, err := s.verifyProjectAccess(ctx, identity, template.ProjectID); err != nil {
		return nil, err
	}

	vars, err := s.repo.GetVariablesByTemplateID(ctx, template.ID)
	if err != nil {
		return nil, errors.Internal("get template variables by template id")
	}
	if vars == nil {
		return nil, errors.NotFound("template variables not found")
	}
	varsResponse := make([]TemplateVariableResponse, len(vars))
	for i, v := range vars {
		varsResponse[i] = v.ToResponse()
	}

	return &TemplateResponse{
		ID:                       template.ID,
		Name:                     template.Name,
		Subject:                  template.Subject,
		BodyHtml:                 template.BodyHtml,
		BodyText:                 template.BodyText,
		Version:                  template.Version,
		TemplateVariableResponse: varsResponse,
		CreatedAt:                template.CreatedAt,
		UpdatedAt:                template.UpdatedAt,
	}, nil
}

func (s *service) verifyProjectAccess(ctx context.Context, identity auth.Identity, projectID uuid.UUID) (*ProjectResponse, error) {
	project, err := s.projectSvc.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if identity.IsStudent() && project.UserID != identity.UserID {
		return nil, errors.Forbidden("not the owner")
	}
	return project, nil
}
