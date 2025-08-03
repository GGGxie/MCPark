package handler

import (
	"net/http"

	"MCPark/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler 处理器结构体，包含所有需要的服务
type Handler struct {
	mcpService *service.MCPService
}

// NewHandler 创建新的处理器实例
func NewHandler(mcpService *service.MCPService) *Handler {
	return &Handler{
		mcpService: mcpService,
	}
}

// Response 标准API响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse 成功响应
func (h *Handler) SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// ErrorResponse 错误响应
func (h *Handler) ErrorResponse(c *gin.Context, code int, message string, err error) {
	response := Response{
		Code:    code,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	var statusCode int
	switch {
	case code >= 400 && code < 500:
		statusCode = code
	case code >= 500:
		statusCode = http.StatusInternalServerError
	default:
		statusCode = http.StatusBadRequest
	}

	c.JSON(statusCode, response)
}
