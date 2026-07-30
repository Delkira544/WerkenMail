package middleware

import (
	sharedauth "github.com/Delkira544/rakiduam/internal/shared/auth"
	"github.com/Delkira544/rakiduam/internal/shared/errors"
	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...sharedauth.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := sharedauth.Role(c.GetString("role"))
		for _, r := range roles {
			if r == userRole {
				c.Next()
				return
			}
		}
		c.Error(errors.Forbidden("insufficient permissions"))
		c.Abort()
	}
}
