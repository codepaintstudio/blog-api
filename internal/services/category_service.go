package services

import (
	"errors"
	"time"

	"blog-api/internal/models"
	"blog-api/internal/repositories"
	"blog-api/pkg/logger"
	
	"gorm.io/gorm"
)

// CategoryService 分类服务接口
type CategoryService interface {
	Create(req *CreateCategoryRequest) (*models.Category, error)
	GetByID(id uint) (*models.Category, error)
	Update(id uint, req *UpdateCategoryRequest) (*models.Category, error)
	Delete(id uint) error
	List() ([]*CategoryListResponse, error)
	ListActive() ([]*CategoryListResponse, error)
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description" validate:"max=500"`
	Color       string `json:"color" validate:"max=7"`
	Icon        string `json:"icon" validate:"max=50"`
	Sort        int    `json:"sort"`
	IsActive    bool   `json:"is_active"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name" validate:"max=50"`
	Description string `json:"description" validate:"max=500"`
	Color       string `json:"color" validate:"max=7"`
	Icon        string `json:"icon" validate:"max=50"`
	Sort        *int   `json:"sort"`
	IsActive    *bool  `json:"is_active"`
}

// CategoryResponse 分类响应
type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	Icon        string    `json:"icon"`
	Sort        int       `json:"sort"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryListResponse 分类列表响应
type CategoryListResponse struct {
	*CategoryResponse
	ArticleCount int64 `json:"article_count"`
}

// categoryService 分类服务实现
type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

// NewCategoryService 创建新的分类服务
func NewCategoryService() CategoryService {
	return &categoryService{
		categoryRepo: repositories.NewCategoryRepository(),
	}
}

// Create 创建分类
func (s *categoryService) Create(req *CreateCategoryRequest) (*models.Category, error) {
	// 检查分类名称是否已存在
	exists, err := s.categoryRepo.ExistsByName(req.Name)
	if err != nil {
		logger.Error("检查分类名称失败", logger.Err("error", err))
		return nil, err
	}
	if exists {
		return nil, errors.New("分类名称已存在")
	}
	
	// 设置默认值
	color := "#409EFF"
	if req.Color != "" {
		color = req.Color
	}
	
	// 创建分类
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		Color:       color,
		Icon:        req.Icon,
		Sort:        req.Sort,
		IsActive:    req.IsActive,
	}
	
	if err := s.categoryRepo.Create(category); err != nil {
		logger.Error("创建分类失败", logger.Err("error", err))
		return nil, err
	}
	
	logger.Info("分类创建成功",
		logger.Int("category_id", int(category.ID)),
		logger.String("name", category.Name),
	)
	
	return category, nil
}

// GetByID 根据ID获取分类
func (s *categoryService) GetByID(id uint) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, err
	}
	return category, nil
}

// Update 更新分类
func (s *categoryService) Update(id uint, req *UpdateCategoryRequest) (*models.Category, error) {
	// 获取分类
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, err
	}
	
	// 检查分类名称是否已存在（排除当前分类）
	if req.Name != "" && req.Name != category.Name {
		exists, err := s.categoryRepo.ExistsByName(req.Name)
		if err != nil {
			logger.Error("检查分类名称失败", logger.Err("error", err))
			return nil, err
		}
		if exists {
			return nil, errors.New("分类名称已存在")
		}
	}
	
	// 更新字段
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Color != "" {
		category.Color = req.Color
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.Sort != nil {
		category.Sort = *req.Sort
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	
	if err := s.categoryRepo.Update(category); err != nil {
		logger.Error("更新分类失败", logger.Err("error", err))
		return nil, err
	}
	
	logger.Info("分类更新成功",
		logger.Int("category_id", int(id)),
		logger.String("name", category.Name),
	)
	
	return category, nil
}

// Delete 删除分类
func (s *categoryService) Delete(id uint) error {
	// 获取分类
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("分类不存在")
		}
		return err
	}
	
	// 检查是否有文章使用该分类
	articleCount, err := s.categoryRepo.GetArticleCount(id)
	if err != nil {
		logger.Error("检查分类文章数量失败", logger.Err("error", err))
		return err
	}
	
	if articleCount > 0 {
		return errors.New("该分类下还有文章，无法删除")
	}
	
	if err := s.categoryRepo.Delete(id); err != nil {
		logger.Error("删除分类失败", logger.Err("error", err))
		return err
	}
	
	logger.Info("分类删除成功",
		logger.Int("category_id", int(id)),
		logger.String("name", category.Name),
	)
	
	return nil
}

// List 获取所有分类列表
func (s *categoryService) List() ([]*CategoryListResponse, error) {
	categoriesWithCount, err := s.categoryRepo.ListWithArticleCount()
	if err != nil {
		return nil, err
	}
	
	responses := make([]*CategoryListResponse, len(categoriesWithCount))
	for i, category := range categoriesWithCount {
		responses[i] = &CategoryListResponse{
			CategoryResponse: &CategoryResponse{
				ID:          category.ID,
				Name:        category.Name,
				Description: category.Description,
				Color:       category.Color,
				Icon:        category.Icon,
				Sort:        category.Sort,
				IsActive:    category.IsActive,
				CreatedAt:   category.CreatedAt,
				UpdatedAt:   category.UpdatedAt,
			},
			ArticleCount: category.ArticleCount,
		}
	}
	
	return responses, nil
}

// ListActive 获取启用的分类列表
func (s *categoryService) ListActive() ([]*CategoryListResponse, error) {
	categories, err := s.categoryRepo.ListActive()
	if err != nil {
		return nil, err
	}
	
	responses := make([]*CategoryListResponse, len(categories))
	for i, category := range categories {
		// 获取文章数量
		articleCount, _ := s.categoryRepo.GetArticleCount(category.ID)
		
		responses[i] = &CategoryListResponse{
			CategoryResponse: &CategoryResponse{
				ID:          category.ID,
				Name:        category.Name,
				Description: category.Description,
				Color:       category.Color,
				Icon:        category.Icon,
				Sort:        category.Sort,
				IsActive:    category.IsActive,
				CreatedAt:   category.CreatedAt,
				UpdatedAt:   category.UpdatedAt,
			},
			ArticleCount: articleCount,
		}
	}
	
	return responses, nil
}