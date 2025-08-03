package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck 健康检查接口
// @Summary 健康检查
// @Description 检查API服务状态
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /api/v1/health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "MCP Info API is running",
	})
}

// ServeIndex 服务首页
func (h *Handler) ServeIndex(c *gin.Context) {
	c.File("./static/index.html")
}

// HandleNotFound 处理前端路由的404
func (h *Handler) HandleNotFound(c *gin.Context) {
	path := c.Request.URL.Path

	// 如果是静态资源请求，返回404
	if len(path) >= 7 && path[:7] == "/assets" {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Static resource not found",
		})
		return
	}

	// 否则返回前端应用
	c.File("./static/index.html")
}
