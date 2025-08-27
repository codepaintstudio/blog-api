package repositories

import (
	"blog-api/internal/models"
	"blog-api/pkg/database"
	"gorm.io/gorm"
)

// CommentRepository 评论仓库接口
type CommentRepository interface {
	Create(comment *models.Comment) error
	GetByID(id uint) (*models.Comment, error)
	Update(comment *models.Comment) error
	Delete(id uint) error
	GetByArticleID(articleID uint, page, size int) ([]*models.Comment, int64, error)
	GetReplies(parentID uint, page, size int) ([]*models.Comment, int64, error)
	GetUserComments(userID uint, page, size int) ([]*models.Comment, int64, error)
}

// commentRepository 评论仓库实现
type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓库实例
func NewCommentRepository() CommentRepository {
	return &commentRepository{
		db: database.GetDB(),
	}
}

// Create 创建评论
func (r *commentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

// GetByID 根据ID获取评论
func (r *commentRepository) GetByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Preload("User").
		Preload("Article").
		Preload("Parent").
		First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// Update 更新评论
func (r *commentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

// Delete 删除评论（软删除）
func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}

// GetByArticleID 获取文章的评论列表
func (r *commentRepository) GetByArticleID(articleID uint, page, size int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	// 计算总数（只统计顶级评论）
	if err := r.db.Model(&models.Comment{}).
		Where("article_id = ? AND parent_id IS NULL", articleID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询顶级评论
	offset := (page - 1) * size
	err := r.db.Where("article_id = ? AND parent_id IS NULL", articleID).
		Preload("User").
		Preload("Children").
		Preload("Children.User").
		Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&comments).Error

	return comments, total, err
}

// GetReplies 获取评论的回复列表
func (r *commentRepository) GetReplies(parentID uint, page, size int) ([]*models.Comment, int64, error) {
	var replies []*models.Comment
	var total int64

	// 计算总数
	if err := r.db.Model(&models.Comment{}).
		Where("parent_id = ?", parentID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询回复
	offset := (page - 1) * size
	err := r.db.Where("parent_id = ?", parentID).
		Preload("User").
		Order("created_at ASC").
		Limit(size).
		Offset(offset).
		Find(&replies).Error

	return replies, total, err
}

// GetUserComments 获取用户的评论列表
func (r *commentRepository) GetUserComments(userID uint, page, size int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	// 计算总数
	if err := r.db.Model(&models.Comment{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	err := r.db.Where("user_id = ?", userID).
		Preload("Article").
		Preload("Parent").
		Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&comments).Error

	return comments, total, err
}