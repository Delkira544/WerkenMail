package middleware

import (
	"strings"

	"github.com/Delkira544/rakiduam/internal/shared/errors"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Error(errors.Unauthorized("missing authorization header"))
			c.Abort()
			return
		}

		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || token == "" {
			c.Error(errors.Unauthorized("invalid authorization header format"))
			c.Abort()
			return
		}

		// TODO: validate JWT token and extract user info
		c.Set("user_id", token)
		c.Next()
	}
}
