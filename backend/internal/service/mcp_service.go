package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"MCPark/internal/model"
	"MCPark/internal/repository"

	"gorm.io/gorm"
)

type MCPService struct {
	serviceRepo  *repository.ServiceRepository
	categoryRepo *repository.CategoryRepository
}

// NewMCPService 创建MCP服务业务层
func NewMCPService(serviceRepo *repository.ServiceRepository, categoryRepo *repository.CategoryRepository) *MCPService {
	return &MCPService{
		serviceRepo:  serviceRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateServiceRequest 创建服务请求结构
type CreateServiceRequest struct {
	Name             string `json:"name" binding:"required"`
	DisplayName      string `json:"display_name"`
	Description      string `json:"description"`
	ConfigJSON       string `json:"config_json" binding:"required"`
	CategoryName     string `json:"category"`
	Tags             string `json:"tags"`
	Author           string `json:"author"`
	RepositoryURL    string `json:"repository_url"`
	DocumentationURL string `json:"documentation_url"`
	Version          string `json:"version"`
}

// CreateService 创建新的MCP服务
func (s *MCPService) CreateService(req *CreateServiceRequest) (*model.MCPService, error) {
	// 验证JSON配置格式
	if !s.isValidJSON(req.ConfigJSON) {
		return nil, errors.New("invalid JSON configuration")
	}

	// 检查服务名称是否已存在
	existing, err := s.serviceRepo.GetByName(req.Name)
	if err == nil && existing != nil {
		return nil, errors.New("service name already exists")
	}

	// 处理分类
	var categoryID *uint
	if req.CategoryName != "" {
		category, err := s.categoryRepo.GetByName(req.CategoryName)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("failed to get category: %v", err)
			}
		} else {
			categoryID = &category.ID
		}
	}

	// 如果没有设置DisplayName，使用Name
	displayName := req.DisplayName
	if displayName == "" {
		displayName = req.Name
	}

	// 创建服务实例
	service := &model.MCPService{
		Name:             req.Name,
		DisplayName:      displayName,
		Description:      req.Description,
		ConfigJSON:       req.ConfigJSON,
		CategoryID:       categoryID,
		Tags:             req.Tags,
		Author:           req.Author,
		RepositoryURL:    req.RepositoryURL,
		DocumentationURL: req.DocumentationURL,
		Version:          req.Version,
		InstallCount:     0,
		IsVerified:       false,
		IsActive:         true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// 保存到数据库
	if err := s.serviceRepo.Create(service); err != nil {
		return nil, fmt.Errorf("failed to create service: %v", err)
	}

	// 返回完整的服务信息（包含分类）
	return s.serviceRepo.GetByID(service.ID)
}

// GetAllServices 获取所有服务
func (s *MCPService) GetAllServices() ([]model.MCPService, error) {
	return s.serviceRepo.GetAll()
}

// GetServiceByID 根据ID获取服务
func (s *MCPService) GetServiceByID(id uint) (*model.MCPService, error) {
	return s.serviceRepo.GetByID(id)
}

// UpdateService 更新服务
func (s *MCPService) UpdateService(id uint, req *CreateServiceRequest) (*model.MCPService, error) {
	// 获取现有服务
	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("service not found: %v", err)
	}

	// 验证JSON配置格式
	if !s.isValidJSON(req.ConfigJSON) {
		return nil, errors.New("invalid JSON configuration")
	}

	// 检查服务名称是否已被其他服务占用
	if req.Name != service.Name {
		existing, err := s.serviceRepo.GetByName(req.Name)
		if err == nil && existing != nil && existing.ID != id {
			return nil, errors.New("service name already exists")
		}
	}

	// 处理分类
	var categoryID *uint
	if req.CategoryName != "" {
		category, err := s.categoryRepo.GetByName(req.CategoryName)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("failed to get category: %v", err)
			}
		} else {
			categoryID = &category.ID
		}
	}

	// 如果没有设置DisplayName，使用Name
	displayName := req.DisplayName
	if displayName == "" {
		displayName = req.Name
	}

	// 更新服务信息
	service.Name = req.Name
	service.DisplayName = displayName
	service.Description = req.Description
	service.ConfigJSON = req.ConfigJSON
	service.CategoryID = categoryID
	service.Tags = req.Tags
	service.Author = req.Author
	service.RepositoryURL = req.RepositoryURL
	service.DocumentationURL = req.DocumentationURL
	service.Version = req.Version
	service.UpdatedAt = time.Now()

	// 保存更新
	if err := s.serviceRepo.Update(service); err != nil {
		return nil, fmt.Errorf("failed to update service: %v", err)
	}

	// 返回更新后的服务信息
	return s.serviceRepo.GetByID(service.ID)
}

// DeleteService 删除服务
func (s *MCPService) DeleteService(id uint) error {
	// 检查服务是否存在
	_, err := s.serviceRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("service not found: %v", err)
	}

	return s.serviceRepo.Delete(id)
}

// SearchServices 搜索服务
func (s *MCPService) SearchServices(keyword string) ([]model.MCPService, error) {
	if keyword == "" {
		return s.GetAllServices()
	}
	return s.serviceRepo.Search(strings.TrimSpace(keyword))
}

// IncrementInstallCount 增加安装次数
func (s *MCPService) IncrementInstallCount(id uint) error {
	// 检查服务是否存在
	_, err := s.serviceRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("service not found: %v", err)
	}

	return s.serviceRepo.IncrementInstallCount(id)
}

// GetAllCategories 获取所有分类
func (s *MCPService) GetAllCategories() ([]model.Category, error) {
	return s.categoryRepo.GetAll()
}

// isValidJSON 验证JSON格式
func (s *MCPService) isValidJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
