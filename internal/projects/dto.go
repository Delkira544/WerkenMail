package projects

import (
	"time"

	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type ListProjectRequest struct {
	Search *string `form:"search"`
	Limit  *int    `form:"limit"`
	Offset *int    `form:"offset"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type ProjectResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toResponse(projects []Project) []ProjectResponse {
	responses := make([]ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = *project.ToResponse()
	}
	return responses
}
