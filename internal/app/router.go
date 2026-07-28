package app

import (
	"net/http"

	"github.com/Delkira544/rakiduam/internal/middleware"
	"github.com/gin-gonic/gin"
)

// NewRouter arma el gin.Engine: middlewares globales, health check y los
// grupos de rutas versionadas de cada feature.
func NewRouter(h *Handlers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorInterceptor())

	r.GET("/health", healthCheck)

	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		authGroup.POST("/login", h.Auth.Login)
	}

	return r
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
