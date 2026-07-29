package app

import (
	"net/http"

	"github.com/Delkira544/rakiduam/config"
	"github.com/Delkira544/rakiduam/internal/middleware"
	"github.com/Delkira544/rakiduam/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// NewRouter arma el gin.Engine: middlewares globales, health check y los
// grupos de rutas versionadas de cada feature.
func NewRouter(h *Handlers, cfg config.Config) *gin.Engine {
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
	prueba := v1.Group("/prueba")
	prueba.Use(middleware.AuthRequired(cfg.App.JWTSecret))
	{
		prueba.GET("/student", middleware.RequireRole("student"), func(c *gin.Context) {
			response.OK(c, "Acceso permitido estudiante")
		})
		prueba.GET("/func", middleware.RequireRole("func"), func(c *gin.Context) {
			response.OK(c, "Acceso permitido funcionario")
		})
	}

	return r
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
