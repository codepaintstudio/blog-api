package repositories

import (
	"blog-api/internal/models"
	"blog-api/pkg/database"
	
	"gorm.io/gorm"
)

// ArticleRepository 文章仓库接口
type ArticleRepository interface {
	Create(article *models.Article) error
	GetByID(id uint) (*models.Article, error)
	GetByIDWithAuthor(id uint) (*models.Article, error)
	Update(article *models.Article) error
	Delete(id uint) error
	List(offset, limit int, filters *ArticleFilters) ([]*models.Article, int64, error)
	GetByUserID(userID uint, offset, limit int) ([]*models.Article, int64, error)
	GetByCategoryID(categoryID uint, offset, limit int) ([]*models.Article, int64, error)
	Search(keyword string, offset, limit int) ([]*models.Article, int64, error)
	IncrementViewCount(id uint) error
	IncrementLikeCount(id uint) error
	DecrementLikeCount(id uint) error
	IncrementFavoriteCount(id uint) error
	DecrementFavoriteCount(id uint) error
	IncrementCommentCount(id uint) error
	DecrementCommentCount(id uint) error
}

// ArticleFilters 文章查询过滤器
type ArticleFilters struct {
	Status     string
	Visibility string
	CategoryID *uint
	UserID     *uint
	IsTop      *bool
	SortBy     string // created_at, updated_at, view_count, like_count
	SortOrder  string // asc, desc
}

// articleRepository 文章仓库实现
type articleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建新的文章仓库
func NewArticleRepository() ArticleRepository {
	return &articleRepository{
		db: database.GetDB(),
	}
}

// Create 创建文章
func (r *articleRepository) Create(article *models.Article) error {
	return r.db.Create(article).Error
}

// GetByID 根据ID获取文章
func (r *articleRepository) GetByID(id uint) (*models.Article, error) {
	var article models.Article
	err := r.db.First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// GetByIDWithAuthor 根据ID获取文章（包含作者信息）
func (r *articleRepository) GetByIDWithAuthor(id uint) (*models.Article, error) {
	var article models.Article
	err := r.db.Preload("User").Preload("Category").First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// Update 更新文章
func (r *articleRepository) Update(article *models.Article) error {
	return r.db.Save(article).Error
}

// Delete 删除文章（软删除）
func (r *articleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Article{}, id).Error
}

// List 获取文章列表
func (r *articleRepository) List(offset, limit int, filters *ArticleFilters) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64
	
	query := r.db.Model(&models.Article{}).Preload("User").Preload("Category")
	
	// 应用过滤器
	if filters != nil {
		if filters.Status != "" {
			query = query.Where("status = ?", filters.Status)
		}
		if filters.Visibility != "" {
			query = query.Where("visibility = ?", filters.Visibility)
		}
		if filters.CategoryID != nil {
			query = query.Where("category_id = ?", *filters.CategoryID)
		}
		if filters.UserID != nil {
			query = query.Where("user_id = ?", *filters.UserID)
		}
		if filters.IsTop != nil {
			query = query.Where("is_top = ?", *filters.IsTop)
		}
		
		// 排序
		if filters.SortBy != "" {
			order := "DESC"
			if filters.SortOrder == "asc" {
				order = "ASC"
			}
			query = query.Order(filters.SortBy + " " + order)
		} else {
			// 默认排序：置顶文章在前，然后按创建时间倒序
			query = query.Order("is_top DESC, created_at DESC")
		}
	} else {
		query = query.Order("is_top DESC, created_at DESC")
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := query.Offset(offset).Limit(limit).Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}
	
	return articles, total, nil
}

// GetByUserID 根据用户ID获取文章列表
func (r *articleRepository) GetByUserID(userID uint, offset, limit int) ([]*models.Article, int64, error) {
	filters := &ArticleFilters{
		UserID: &userID,
	}
	return r.List(offset, limit, filters)
}

// GetByCategoryID 根据分类ID获取文章列表
func (r *articleRepository) GetByCategoryID(categoryID uint, offset, limit int) ([]*models.Article, int64, error) {
	filters := &ArticleFilters{
		CategoryID: &categoryID,
		Status:     "published",
		Visibility: "public",
	}
	return r.List(offset, limit, filters)
}

// Search 搜索文章
func (r *articleRepository) Search(keyword string, offset, limit int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64
	
	query := r.db.Model(&models.Article{}).
		Preload("User").
		Preload("Category").
		Where("status = ? AND visibility = ?", "published", "public").
		Where("title LIKE ? OR content LIKE ? OR description LIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Order("is_top DESC, created_at DESC")
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := query.Offset(offset).Limit(limit).Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}
	
	return articles, total, nil
}

// IncrementViewCount 增加浏览量
func (r *articleRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// IncrementLikeCount 增加点赞数
func (r *articleRepository) IncrementLikeCount(id uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", id).Update("like_count", gorm.Expr("like_count + ?", 1)).Error
}

// DecrementLikeCount 减少点赞数
func (r *articleRepository) DecrementLikeCount(id uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", id).Update("like_count", gorm.Expr("like_count - ?", 1)).Error
}

// IncrementFavoriteCount 增加收藏数
func (r *articleRepository) IncrementFavoriteCount(id uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", id).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
}

// DecrementFavoriteCount 减少收藏数
func (r *articleRepository) DecrementFavoriteCount(id uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", id).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
}

// IncrementCommentCount 增加评论数
func (r *articleRepository) IncrementCommentCount(id uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", id).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
}

// DecrementCommentCount 减少评论数
func (r *articleRepository) DecrementCommentCount(id uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", id).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
}