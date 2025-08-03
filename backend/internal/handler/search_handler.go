package handler

import (
	"github.com/gin-gonic/gin"
)

// SearchServices 搜索服务
// @Summary 搜索服务
// @Description 根据关键词搜索MCP服务
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "搜索关键词"
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/search [get]
func (h *Handler) SearchServices(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		h.ErrorResponse(c, 400, "搜索关键词不能为空", nil)
		return
	}

	services, err := h.mcpService.SearchServices(keyword)
	if err != nil {
		h.ErrorResponse(c, 500, "搜索失败", err)
		return
	}

	h.SuccessResponse(c, services)
}
