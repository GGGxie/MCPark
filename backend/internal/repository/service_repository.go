package repository

import (
	"MCPark/internal/model"

	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

// NewServiceRepository 创建服务仓库
func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// Create 创建新的MCP服务
func (r *ServiceRepository) Create(service *model.MCPService) error {
	return r.db.Create(service).Error
}

// GetByID 根据ID获取MCP服务
func (r *ServiceRepository) GetByID(id uint) (*model.MCPService, error) {
	var service model.MCPService
	err := r.db.Preload("Category").First(&service, id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

// GetAll 获取所有MCP服务
func (r *ServiceRepository) GetAll() ([]model.MCPService, error) {
	var services []model.MCPService
	err := r.db.Preload("Category").Order("created_at DESC").Find(&services).Error
	return services, err
}

// Update 更新MCP服务
func (r *ServiceRepository) Update(service *model.MCPService) error {
	return r.db.Save(service).Error
}

// Delete 删除MCP服务
func (r *ServiceRepository) Delete(id uint) error {
	return r.db.Delete(&model.MCPService{}, id).Error
}

// GetByName 根据名称获取MCP服务
func (r *ServiceRepository) GetByName(name string) (*model.MCPService, error) {
	var service model.MCPService
	err := r.db.Preload("Category").Where("name = ?", name).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

// GetByCategory 根据分类获取MCP服务
func (r *ServiceRepository) GetByCategory(categoryID uint) ([]model.MCPService, error) {
	var services []model.MCPService
	err := r.db.Preload("Category").Where("category_id = ?", categoryID).Order("created_at DESC").Find(&services).Error
	return services, err
}

// Search 搜索MCP服务
func (r *ServiceRepository) Search(keyword string) ([]model.MCPService, error) {
	var services []model.MCPService
	err := r.db.Preload("Category").Where(
		"name LIKE ? OR display_name LIKE ? OR description LIKE ? OR tags LIKE ? OR author LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
	).Order("created_at DESC").Find(&services).Error
	return services, err
}

// IncrementInstallCount 增加安装次数
func (r *ServiceRepository) IncrementInstallCount(id uint) error {
	return r.db.Model(&model.MCPService{}).Where("id = ?", id).UpdateColumn("install_count", gorm.Expr("install_count + 1")).Error
}
