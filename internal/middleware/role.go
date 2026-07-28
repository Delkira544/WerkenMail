package middleware

import (
	"github.com/Delkira544/rakiduam/internal/shared/errors"
	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")
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
