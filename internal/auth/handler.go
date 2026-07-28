package auth

import (
	apperrors "github.com/Delkira544/rakiduam/internal/shared/errors"
	"github.com/Delkira544/rakiduam/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperrors.BadRequest(err.Error()))
		return
	}

	session, err := h.svc.Login(req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, session)
}
