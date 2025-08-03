package handler

import (
	"strconv"

	"MCPark/internal/service"

	"github.com/gin-gonic/gin"
)

// GetServices 获取服务列表
// @Summary 获取服务列表
// @Description 获取所有MCP服务列表
// @Tags services
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/services [get]
func (h *Handler) GetServices(c *gin.Context) {
	services, err := h.mcpService.GetAllServices()
	if err != nil {
		h.ErrorResponse(c, 500, "获取服务列表失败", err)
		return
	}

	h.SuccessResponse(c, services)
}

// CreateService 创建新服务
// @Summary 创建新服务
// @Description 创建一个新的MCP服务
// @Tags services
// @Accept json
// @Produce json
// @Param service body service.CreateServiceRequest true "服务信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/v1/services [post]
func (h *Handler) CreateService(c *gin.Context) {
	var req service.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, 400, "请求参数错误", err)
		return
	}

	createdService, err := h.mcpService.CreateService(&req)
	if err != nil {
		h.ErrorResponse(c, 400, "创建服务失败", err)
		return
	}

	h.SuccessResponse(c, createdService)
}

// GetServiceByID 根据ID获取服务
// @Summary 获取服务详情
// @Description 根据ID获取MCP服务详情
// @Tags services
// @Accept json
// @Produce json
// @Param id path int true "服务ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Router /api/v1/services/{id} [get]
func (h *Handler) GetServiceByID(c *gin.Context) {
	id, err := h.parseIDParam(c)
	if err != nil {
		h.ErrorResponse(c, 400, "无效的服务ID", err)
		return
	}

	service, err := h.mcpService.GetServiceByID(id)
	if err != nil {
		h.ErrorResponse(c, 404, "服务不存在", err)
		return
	}

	h.SuccessResponse(c, service)
}

// UpdateService 更新服务
// @Summary 更新服务
// @Description 更新MCP服务信息
// @Tags services
// @Accept json
// @Produce json
// @Param id path int true "服务ID"
// @Param service body service.CreateServiceRequest true "服务信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/v1/services/{id} [put]
func (h *Handler) UpdateService(c *gin.Context) {
	id, err := h.parseIDParam(c)
	if err != nil {
		h.ErrorResponse(c, 400, "无效的服务ID", err)
		return
	}

	var req service.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, 400, "请求参数错误", err)
		return
	}

	updatedService, err := h.mcpService.UpdateService(id, &req)
	if err != nil {
		h.ErrorResponse(c, 400, "更新服务失败", err)
		return
	}

	h.SuccessResponse(c, updatedService)
}

// DeleteService 删除服务
// @Summary 删除服务
// @Description 删除MCP服务
// @Tags services
// @Accept json
// @Produce json
// @Param id path int true "服务ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/v1/services/{id} [delete]
func (h *Handler) DeleteService(c *gin.Context) {
	id, err := h.parseIDParam(c)
	if err != nil {
		h.ErrorResponse(c, 400, "无效的服务ID", err)
		return
	}

	err = h.mcpService.DeleteService(id)
	if err != nil {
		h.ErrorResponse(c, 400, "删除服务失败", err)
		return
	}

	h.SuccessResponse(c, gin.H{"message": "服务删除成功"})
}

// IncrementInstallCount 增加安装次数
// @Summary 增加安装次数
// @Description 记录MCP服务安装，增加安装计数
// @Tags services
// @Accept json
// @Produce json
// @Param id path int true "服务ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/v1/services/{id}/install [post]
func (h *Handler) IncrementInstallCount(c *gin.Context) {
	id, err := h.parseIDParam(c)
	if err != nil {
		h.ErrorResponse(c, 400, "无效的服务ID", err)
		return
	}

	err = h.mcpService.IncrementInstallCount(id)
	if err != nil {
		h.ErrorResponse(c, 400, "更新安装次数失败", err)
		return
	}

	h.SuccessResponse(c, gin.H{"message": "安装次数更新成功"})
}

// parseIDParam 解析路径参数中的ID
func (h *Handler) parseIDParam(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
