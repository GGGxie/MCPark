package handler

import (
	"github.com/gin-gonic/gin"
)

// GetCategories 获取分类列表
// @Summary 获取分类列表
// @Description 获取所有MCP服务分类
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/categories [get]
func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.mcpService.GetAllCategories()
	if err != nil {
		h.ErrorResponse(c, 500, "获取分类列表失败", err)
		return
	}

	h.SuccessResponse(c, categories)
}
