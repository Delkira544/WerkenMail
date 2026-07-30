package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "github.com/Delkira544/rakiduam/internal/shared/errors"
)

type Identity struct {
	UserID uuid.UUID
	Role   Role
}

func (i Identity) IsStudent() bool           { return i.Role.IsStudent() }
func (i Identity) IsFunc() bool              { return i.Role.IsFunc() }
func (i Identity) IsOwner(id uuid.UUID) bool { return i.UserID == id }

// FromGin extrae la Identity seteada por middleware.AuthRequired.
func FromGin(c *gin.Context) (Identity, error) {
	idStr := c.GetString("user_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return Identity{}, apperrors.Unauthorized("invalid identity in context")
	}
	return Identity{UserID: id, Role: Role(c.GetString("role"))}, nil
}
