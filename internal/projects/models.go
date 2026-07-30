package projects

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (p *Project) ToResponse() *ProjectResponse {
	return &ProjectResponse{
		ID:          p.ID,
		UserID:      p.UserID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type CreateProjectInput struct {
	UserID      uuid.UUID
	Name        string
	Description string
}

type ProjectFilter struct {
	UserID *string
	Search *string
	Limit  *int
	Offset *int
}
