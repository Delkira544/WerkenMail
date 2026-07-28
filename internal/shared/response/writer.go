package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, NewSuccessResponse(data))
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, NewSuccessResponse(data))
}

func OKWithMeta(c *gin.Context, data any, meta *Meta) {
	c.JSON(http.StatusOK, NewSuccessResponseWithMeta(data, meta))
}
