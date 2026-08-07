package template

import (
	"time"

	"github.com/google/uuid"
)

type Template struct {
	ID        uuid.UUID `db:"id"`
	ProjectID uuid.UUID `db:"project_id"`
	Name      string    `db:"name"`
	Subject   string    `db:"subject"`
	BodyHtml  *string   `db:"body_html"`
	BodyText  *string   `db:"body_text"`
	Version   uint16    `db:"version"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (t *Template) ToListItemResponse(varCount int) TemplateListItemResponse {
	return TemplateListItemResponse{
		ID:            t.ID,
		Name:          t.Name,
		Subject:       t.Subject,
		Version:       t.Version,
		VariableCount: varCount,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
}

type TemplateVariable struct {
	ID           uuid.UUID `db:"id"`
	TemplateID   uuid.UUID `db:"template_id"`
	Key          string    `db:"key"`
	Type         string    `db:"type"`
	DefaultValue *string   `db:"default_value"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (v *TemplateVariable) ToResponse() TemplateVariableResponse {
	return TemplateVariableResponse{
		ID:           v.ID,
		Key:          v.Key,
		Type:         v.Type,
		DefaultValue: v.DefaultValue,
	}
}

type templateVariableCount struct {
	TemplateID uuid.UUID `db:"template_id"`
	Count      int       `db:"count"`
}
