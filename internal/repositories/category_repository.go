package repositories

import (
	"blog-api/internal/models"
	"blog-api/pkg/database"
	
	"gorm.io/gorm"
)

// CategoryRepository 分类仓库接口
type CategoryRepository interface {
	Create(category *models.Category) error
	GetByID(id uint) (*models.Category, error)
	GetByName(name string) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
	List() ([]*models.Category, error)
	ListActive() ([]*models.Category, error)
	ListWithArticleCount() ([]*CategoryWithCount, error)
	ExistsByName(name string) (bool, error)
	GetArticleCount(categoryID uint) (int64, error)
}

// CategoryWithCount 带文章数量的分类
type CategoryWithCount struct {
	*models.Category
	ArticleCount int64 `json:"article_count"`
}

// categoryRepository 分类仓库实现
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建新的分类仓库
func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{
		db: database.GetDB(),
	}
}

// Create 创建分类
func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// GetByID 根据ID获取分类
func (r *categoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetByName 根据名称获取分类
func (r *categoryRepository) GetByName(name string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Update 更新分类
func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete 删除分类（软删除）
func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

// List 获取所有分类列表
func (r *categoryRepository) List() ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.Order("sort ASC, created_at DESC").Find(&categories).Error
	return categories, err
}

// ListActive 获取启用的分类列表
func (r *categoryRepository) ListActive() ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.Where("is_active = ?", true).Order("sort ASC, created_at DESC").Find(&categories).Error
	return categories, err
}

// ListWithArticleCount 获取带文章数量的分类列表
func (r *categoryRepository) ListWithArticleCount() ([]*CategoryWithCount, error) {
	var results []*CategoryWithCount
	
	err := r.db.Table("categories").
		Select("categories.*, COUNT(articles.id) as article_count").
		Joins("LEFT JOIN articles ON categories.id = articles.category_id AND articles.deleted_at IS NULL AND articles.status = 'published' AND articles.visibility = 'public'").
		Where("categories.deleted_at IS NULL").
		Group("categories.id").
		Order("categories.sort ASC, categories.created_at DESC").
		Scan(&results).Error
	
	return results, err
}

// ExistsByName 检查分类名称是否已存在
func (r *categoryRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// GetArticleCount 获取分类下的文章数量
func (r *categoryRepository) GetArticleCount(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Article{}).
		Where("category_id = ? AND status = ? AND visibility = ?", categoryID, "published", "public").
		Count(&count).Error
	return count, err
}