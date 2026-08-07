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

	template, err := h.service.Create(c.Request.Context(), identity, id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, template)
}

func (h *Handler) ListTemplatesByProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.Error(errors.BadRequest("Invalid project ID"))
		return
	}

	identity, err := auth.FromGin(c)
	if err != nil {
		c.Error(err)
		return
	}

	templates, err := h.service.ListByProjectID(c.Request.Context(), identity, id)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, templates)
}

func (h *Handler) GetTemplate(c *gin.Context) {
	templateID, err := uuid.Parse(c.Param("template_id"))
	if err != nil {
		c.Error(errors.BadRequest("Invalid template ID"))
		return
	}

	identity, err := auth.FromGin(c)
	if err != nil {
		c.Error(err)
		return
	}

	template, err := h.service.GetByID(c.Request.Context(), identity, templateID)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, template)
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup, authMidl gin.HandlerFunc) {
	tempProjects := r.Group("/projects/:project_id/templates")
	tempProjects.Use(authMidl)
	{
		tempProjects.GET("", h.ListTemplatesByProject)
		tempProjects.POST("", h.CreateTemplate)
	}

	templates := r.Group("/templates")
	templates.Use(authMidl)
	{
		templates.GET("/:template_id", h.GetTemplate)
	}

}
