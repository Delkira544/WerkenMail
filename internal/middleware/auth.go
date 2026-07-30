package middleware

import (
	"fmt"
	"strings"

	sharedauth "github.com/Delkira544/rakiduam/internal/shared/auth"
	"github.com/Delkira544/rakiduam/internal/shared/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Error(errors.Unauthorized("missing authorization header"))
			c.Abort()
			return
		}
		tokenString, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || tokenString == "" {
			c.Error(errors.Unauthorized("invalid authorization header format"))
			c.Abort()
			return
		}
		claims := &sharedauth.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.Error(errors.Unauthorized("invalid or expired token"))
			c.Abort()
			return
		}
		c.Set("user_id", claims.Subject)
		c.Set("role", claims.Role.String())
		c.Next()
	}
}
