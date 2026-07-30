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

type TemplateVariable struct {
	ID            uuid.UUID `db:"id"`
	TemplateID    uuid.UUID `db:"template_id"`
	Key           string    `db:"key"`
	Type          string    `db:"type"`
	Required      bool      `db:"required"`
	Default_Value *string   `db:"default_value"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
