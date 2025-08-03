package repository

import (
	"MCPark/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓库
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create 创建新的分类
func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

// GetByID 根据ID获取分类
func (r *CategoryRepository) GetByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAll 获取所有分类
func (r *CategoryRepository) GetAll() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Order("sort_order ASC, created_at ASC").Find(&categories).Error
	return categories, err
}

// Update 更新分类
func (r *CategoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

// Delete 删除分类
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

// GetByName 根据名称获取分类
func (r *CategoryRepository) GetByName(name string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}
