package app

import (
	"net/http"

	"github.com/Delkira544/rakiduam/config"
	"github.com/Delkira544/rakiduam/internal/middleware"
	"github.com/Delkira544/rakiduam/internal/platform/logger"
	"github.com/Delkira544/rakiduam/internal/shared/auth"
	"github.com/Delkira544/rakiduam/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// NewRouter arma el gin.Engine: middlewares globales, health check y los
// grupos de rutas versionadas de cada feature.
func NewRouter(h *Handlers, cfg config.Config) *gin.Engine {
	r := gin.New()
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(logger.L()))
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorInterceptor())

	r.GET("/health", healthCheck)

	v1 := r.Group("/api/v1")
	h.Auth.RegisterRoutes(v1)
	h.Project.RegisterRoutes(v1, middleware.AuthRequired(cfg.App.JWTSecret))
	h.Template.RegisterRoutes(v1, middleware.AuthRequired(cfg.App.JWTSecret))

	prueba := v1.Group("/prueba")
	prueba.Use(middleware.AuthRequired(cfg.App.JWTSecret))
	{
		prueba.GET("/student", middleware.RequireRole(auth.RoleStudent), func(c *gin.Context) {
			response.OK(c, "Acceso permitido estudiante")
		})
		prueba.GET("/func", middleware.RequireRole(auth.RoleFunc), func(c *gin.Context) {
			response.OK(c, "Acceso permitido funcionario")
		})
	}

	return r
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
