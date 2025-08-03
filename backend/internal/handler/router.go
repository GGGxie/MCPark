package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func (h *Handler) SetupRoutes() *gin.Engine {
	r := gin.Default()

	// 静态文件服务 (前端构建产物)
	r.Static("/assets", "./static/assets")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	// 首页路由 - 返回前端应用
	r.GET("/", h.ServeIndex)

	// API健康检查
	r.GET("/api/v1/health", h.HealthCheck)

	// API路由组
	api := r.Group("/api/v1")
	{
		// 服务管理API
		api.GET("/services", h.GetServices)
		api.POST("/services", h.CreateService)
		api.GET("/services/:id", h.GetServiceByID)
		api.PUT("/services/:id", h.UpdateService)
		api.DELETE("/services/:id", h.DeleteService)

		// 安装计数API
		api.POST("/services/:id/install", h.IncrementInstallCount)

		// 分类管理API
		api.GET("/categories", h.GetCategories)

		// 搜索API
		api.GET("/search", h.SearchServices)
	}

	// 处理前端路由 (SPA)
	r.NoRoute(h.HandleNotFound)

	return r
}

// HandleAPINotFound 处理API 404
func (h *Handler) HandleAPINotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"code":    404,
		"message": "API not found",
	})
}
