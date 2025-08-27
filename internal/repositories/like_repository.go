package repositories

import (
	"blog-api/internal/models"
	"blog-api/pkg/database"
	"gorm.io/gorm"
)

// ArticleLikeRepository 文章点赞仓库接口
type ArticleLikeRepository interface {
	Create(like *models.ArticleLike) error
	Delete(userID, articleID uint) error
	GetByUserAndArticle(userID, articleID uint) (*models.ArticleLike, error)
	GetLikeCount(articleID uint) (int64, error)
	GetUserLikes(userID uint, page, size int) ([]*models.ArticleLike, int64, error)
	IsLikedByUser(userID, articleID uint) (bool, error)
	GetArticleLikers(articleID uint, page, size int) ([]*models.User, int64, error)
}

// articleLikeRepository 文章点赞仓库实现
type articleLikeRepository struct {
	db *gorm.DB
}

// NewArticleLikeRepository 创建文章点赞仓库实例
func NewArticleLikeRepository() ArticleLikeRepository {
	return &articleLikeRepository{
		db: database.GetDB(),
	}
}

// Create 创建点赞记录
func (r *articleLikeRepository) Create(like *models.ArticleLike) error {
	return r.db.Create(like).Error
}

// Delete 删除点赞记录
func (r *articleLikeRepository) Delete(userID, articleID uint) error {
	return r.db.Where("user_id = ? AND article_id = ?", userID, articleID).
		Delete(&models.ArticleLike{}).Error
}

// GetByUserAndArticle 获取用户对文章的点赞记录
func (r *articleLikeRepository) GetByUserAndArticle(userID, articleID uint) (*models.ArticleLike, error) {
	var like models.ArticleLike
	err := r.db.Where("user_id = ? AND article_id = ?", userID, articleID).
		First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// GetLikeCount 获取文章点赞数
func (r *articleLikeRepository) GetLikeCount(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.ArticleLike{}).
		Where("article_id = ?", articleID).
		Count(&count).Error
	return count, err
}

// GetUserLikes 获取用户点赞的文章列表
func (r *articleLikeRepository) GetUserLikes(userID uint, page, size int) ([]*models.ArticleLike, int64, error) {
	var likes []*models.ArticleLike
	var total int64

	// 计算总数
	if err := r.db.Model(&models.ArticleLike{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，预加载文章信息
	offset := (page - 1) * size
	err := r.db.Where("user_id = ?", userID).
		Preload("Article").
		Preload("Article.User").
		Preload("Article.Category").
		Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&likes).Error

	return likes, total, err
}

// IsLikedByUser 检查用户是否已点赞文章
func (r *articleLikeRepository) IsLikedByUser(userID, articleID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.ArticleLike{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count).Error
	return count > 0, err
}

// GetArticleLikers 获取点赞文章的用户列表
func (r *articleLikeRepository) GetArticleLikers(articleID uint, page, size int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	// 计算总数
	if err := r.db.Table("article_likes").
		Where("article_id = ?", articleID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询点赞用户
	offset := (page - 1) * size
	err := r.db.Table("users").
		Select("users.*").
		Joins("INNER JOIN article_likes ON users.id = article_likes.user_id").
		Where("article_likes.article_id = ?", articleID).
		Order("article_likes.created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

// CommentLikeRepository 评论点赞仓库接口
type CommentLikeRepository interface {
	Create(like *models.CommentLike) error
	Delete(userID, commentID uint) error
	GetByUserAndComment(userID, commentID uint) (*models.CommentLike, error)
	GetLikeCount(commentID uint) (int64, error)
	GetUserCommentLikes(userID uint, page, size int) ([]*models.CommentLike, int64, error)
	IsLikedByUser(userID, commentID uint) (bool, error)
}

// commentLikeRepository 评论点赞仓库实现
type commentLikeRepository struct {
	db *gorm.DB
}

// NewCommentLikeRepository 创建评论点赞仓库实例
func NewCommentLikeRepository() CommentLikeRepository {
	return &commentLikeRepository{
		db: database.GetDB(),
	}
}

// Create 创建评论点赞记录
func (r *commentLikeRepository) Create(like *models.CommentLike) error {
	return r.db.Create(like).Error
}

// Delete 删除评论点赞记录
func (r *commentLikeRepository) Delete(userID, commentID uint) error {
	return r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).
		Delete(&models.CommentLike{}).Error
}

// GetByUserAndComment 获取用户对评论的点赞记录
func (r *commentLikeRepository) GetByUserAndComment(userID, commentID uint) (*models.CommentLike, error) {
	var like models.CommentLike
	err := r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).
		First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// GetLikeCount 获取评论点赞数
func (r *commentLikeRepository) GetLikeCount(commentID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.CommentLike{}).
		Where("comment_id = ?", commentID).
		Count(&count).Error
	return count, err
}

// GetUserCommentLikes 获取用户点赞的评论列表
func (r *commentLikeRepository) GetUserCommentLikes(userID uint, page, size int) ([]*models.CommentLike, int64, error) {
	var likes []*models.CommentLike
	var total int64

	// 计算总数
	if err := r.db.Model(&models.CommentLike{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，预加载评论信息
	offset := (page - 1) * size
	err := r.db.Where("user_id = ?", userID).
		Preload("Comment").
		Preload("Comment.User").
		Preload("Comment.Article").
		Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&likes).Error

	return likes, total, err
}

// IsLikedByUser 检查用户是否已点赞评论
func (r *commentLikeRepository) IsLikedByUser(userID, commentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.CommentLike{}).
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		Count(&count).Error
	return count > 0, err
}