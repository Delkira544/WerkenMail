package template

import (
	"github.com/Delkira544/rakiduam/internal/shared/auth"
	"github.com/Delkira544/rakiduam/internal/shared/errors"
	"github.com/Delkira544/rakiduam/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.Error(errors.BadRequest("Invalid project ID"))
		return
	}

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest(err.Error()))
		return
	}

	identity, err := auth.FromGin(c)
	if err != nil {
		c.Error(err)
		return
	}

	template, err := h.service.CreateTemplate(c.Request.Context(), identity, id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, template)
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup, authMidl gin.HandlerFunc) {
	template := r.Group("/projects/:project_id/templates")
	template.Use(authMidl)
	{
		template.POST("", h.CreateTemplate)
	}

}
