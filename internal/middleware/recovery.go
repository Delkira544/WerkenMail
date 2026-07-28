package middleware

import (
	"errors"
	"net/http"

	apperrors "github.com/Delkira544/rakiduam/internal/shared/errors"
	"github.com/Delkira544/rakiduam/internal/shared/response"
	"github.com/gin-gonic/gin"
)

func ErrorInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var appErr *apperrors.AppError
			if errors.As(err, &appErr) {
				c.AbortWithStatusJSON(appErr.Status, response.NewErrorResponse(appErr))
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				response.NewErrorResponse(apperrors.Internal("")))
		}
	}
}

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			response.NewErrorResponse(apperrors.Internal("")))
	})
}
