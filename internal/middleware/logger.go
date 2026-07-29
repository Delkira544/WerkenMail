package middleware

import (
	"github.com/Delkira544/rakiduam/internal/platform/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetString("request_id")
		l := log.With(
			zap.String("request_id", rid),
		)
		ctx := logger.WithContext(c.Request.Context(), l)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
