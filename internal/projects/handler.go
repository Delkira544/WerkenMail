package projects

import (
	"github.com/Delkira544/rakiduam/internal/shared/auth"
	"github.com/Delkira544/rakiduam/internal/shared/errors"
	"github.com/Delkira544/rakiduam/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest(err.Error()))
		return
	}

	identity, err := auth.FromGin(c)
	if err != nil {
		c.Error(err)
		return
	}

	input := &CreateProjectInput{
		UserID:      identity.UserID,
		Name:        req.Name,
		Description: req.Description,
	}

	project, err := h.svc.Create(c.Request.Context(), input)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, project)
}

func (h *Handler) ListProjects(c *gin.Context) {
	var req ListProjectRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(errors.BadRequest(err.Error()))
		return
	}

	identity, err := auth.FromGin(c)
	if err != nil {
		c.Error(err)
		return
	}

	projects, err := h.svc.List(c.Request.Context(), identity, req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, projects)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(errors.BadRequest("invalid project ID"))
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest(err.Error()))
		return
	}

	identity, err := auth.FromGin(c)
	if err != nil {
		c.Error(err)
		return
	}

	project, err := h.svc.Update(c.Request.Context(), id, identity, req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, project)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(errors.BadRequest("invalid project ID"))
		return
	}

	identity, err := auth.FromGin(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.svc.Delete(c.Request.Context(), id, identity)
	if err != nil {
		c.Error(err)
		return
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	project := rg.Group("/projects")
	project.Use(authMiddleware)
	{
		project.POST("", h.CreateProject)
		project.GET("", h.ListProjects)
		project.PATCH("/:id", h.UpdateProject)
		project.DELETE("/:id", h.DeleteProject)
	}
}
