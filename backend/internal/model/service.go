package model

import (
	"time"
)

// MCPService MCP服务模型
type MCPService struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Name             string    `json:"name" gorm:"uniqueIndex;not null;size:100" binding:"required"`
	DisplayName      string    `json:"display_name" gorm:"not null;size:200" binding:"required"`
	Description      string    `json:"description" gorm:"type:text"`
	ConfigJSON       string    `json:"config_json" gorm:"type:text;not null" binding:"required"`
	CategoryID       *uint     `json:"category_id"`
	Category         *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Tags             string    `json:"tags" gorm:"size:500"`
	Author           string    `json:"author" gorm:"size:100"`
	RepositoryURL    string    `json:"repository_url" gorm:"size:500"`
	DocumentationURL string    `json:"documentation_url" gorm:"size:500"`
	Version          string    `json:"version" gorm:"size:50"`
	InstallCount     int       `json:"install_count" gorm:"default:0"`
	IsVerified       bool      `json:"is_verified" gorm:"default:false"`
	IsActive         bool      `json:"is_active" gorm:"default:true"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Category 分类模型
type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;size:100" binding:"required"`
	DisplayName string    `json:"display_name" gorm:"not null;size:200" binding:"required"`
	Description string    `json:"description" gorm:"type:text"`
	Icon        string    `json:"icon" gorm:"size:100"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
}

// 移除用户相关模型，私人系统不需要用户管理

// ServiceListRequest 服务列表请求参数
type ServiceListRequest struct {
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=100"`
	CategoryID *uint  `form:"category_id"`
	Tags       string `form:"tags"`
	Search     string `form:"search"`
	SortBy     string `form:"sort_by,default=created_at"` // created_at, install_count, name
	SortOrder  string `form:"sort_order,default=desc"`    // asc, desc
}

// ServiceCreateRequest 创建服务请求
type ServiceCreateRequest struct {
	Name             string `json:"name" binding:"required,max=100"`
	DisplayName      string `json:"display_name" binding:"required,max=200"`
	Description      string `json:"description"`
	ConfigJSON       string `json:"config_json" binding:"required"`
	CategoryID       *uint  `json:"category_id"`
	Tags             string `json:"tags"`
	Author           string `json:"author"`
	RepositoryURL    string `json:"repository_url"`
	DocumentationURL string `json:"documentation_url"`
	Version          string `json:"version"`
}

// ServiceUpdateRequest 更新服务请求
type ServiceUpdateRequest struct {
	DisplayName      *string `json:"display_name"`
	Description      *string `json:"description"`
	ConfigJSON       *string `json:"config_json"`
	CategoryID       *uint   `json:"category_id"`
	Tags             *string `json:"tags"`
	Author           *string `json:"author"`
	RepositoryURL    *string `json:"repository_url"`
	DocumentationURL *string `json:"documentation_url"`
	Version          *string `json:"version"`
	IsActive         *bool   `json:"is_active"`
}

// DeepLinkResponse Deep Link响应
type DeepLinkResponse struct {
	ServiceName string `json:"service_name"`
	DeepLink    string `json:"deep_link"`
	ConfigJSON  string `json:"config_json"`
}

// APIResponse 通用API响应
type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// PagedResponse 分页响应
type PagedResponse struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// TableName 指定表名
func (MCPService) TableName() string {
	return "mcp_services"
}

func (Category) TableName() string {
	return "categories"
}
