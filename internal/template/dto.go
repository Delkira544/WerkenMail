package template

import (
	"time"

	"github.com/Delkira544/rakiduam/internal/shared/emailtemplate"
	"github.com/google/uuid"
)

type CreateTemplateRequest struct {
	Name      string                          `json:"name" validate:"required"`
	Subject   string                          `json:"subject" validate:"required"`
	BodyHtml  *string                         `json:"body_html" validate:"required"`
	BodyText  *string                         `json:"body_text"`
	Variables []CreateTemplateVariableRequest `json:"variables"`
}

type CreateTemplateVariableRequest struct {
	Key          string  `json:"key" validate:"required"`
	Type         string  `json:"type" validate:"required"`
	DefaultValue *string `json:"default_value,omitempty"`
}

func toTemplateVariableRequests(vars []CreateTemplateVariableRequest) ([]emailtemplate.Variable, error) {
	var templateVars []emailtemplate.Variable
	for _, v := range vars {
		tv, err := v.ToTemplateVariable()
		if err != nil {
			return nil, err
		}
		templateVars = append(templateVars, *tv)
	}
	return templateVars, nil
}

func (r *CreateTemplateVariableRequest) ToTemplateVariable() (*emailtemplate.Variable, error) {
	return &emailtemplate.Variable{
		Key:          r.Key,
		Type:         emailtemplate.VariableType(r.Type),
		DefaultValue: r.DefaultValue,
	}, nil
}

type TemplateResponse struct {
	ID                       uuid.UUID                  `json:"id"`
	Name                     string                     `json:"name"`
	Subject                  string                     `json:"subject"`
	BodyHtml                 *string                    `json:"body_html,omitempty"`
	BodyText                 *string                    `json:"body_text,omitempty"`
	Version                  uint16                     `json:"version"`
	TemplateVariableResponse []TemplateVariableResponse `json:"variables"`
	CreatedAt                time.Time                  `json:"created_at"`
	UpdatedAt                time.Time                  `json:"updated_at"`
}

type TemplateListItemResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Subject       string    `json:"subject"`
	VariableCount int       `json:"variable_count"`
	Version       uint16    `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type TemplateVariableResponse struct {
	ID           uuid.UUID `json:"id"`
	Key          string    `json:"key"`
	Type         string    `json:"type"`
	DefaultValue *string   `json:"default_value,omitempty"`
}

// ProjectResponse represents the response structure for a project.
type ProjectResponse struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
