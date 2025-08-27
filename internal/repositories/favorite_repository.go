package repositories

import (
	"blog-api/internal/models"
	"blog-api/pkg/database"
	"gorm.io/gorm"
)

// FavoriteRepository 收藏仓库接口
type FavoriteRepository interface {
	// 收藏夹管理
	CreateFolder(folder *models.FavoriteFolder) error
	GetFolderByID(id uint) (*models.FavoriteFolder, error)
	GetUserFolders(userID uint) ([]*models.FavoriteFolder, error)
	UpdateFolder(folder *models.FavoriteFolder) error
	DeleteFolder(id uint) error
	GetDefaultFolder(userID uint) (*models.FavoriteFolder, error)
	FolderExists(userID uint, name string) (bool, error)
	
	// 文章收藏管理
	AddFavorite(favorite *models.ArticleFavorite) error
	RemoveFavorite(userID, articleID uint) error
	GetFavorite(userID, articleID uint) (*models.ArticleFavorite, error)
	IsFavorited(userID, articleID uint) (bool, error)
	GetFolderFavorites(folderID uint, page, size int) ([]*models.ArticleFavorite, int64, error)
	GetUserFavorites(userID uint, page, size int) ([]*models.ArticleFavorite, int64, error)
	MoveFavorite(userID, articleID, newFolderID uint) error
	GetFavoriteCount(articleID uint) (int64, error)
}

// favoriteRepository 收藏仓库实现
type favoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository 创建收藏仓库实例
func NewFavoriteRepository() FavoriteRepository {
	return &favoriteRepository{
		db: database.GetDB(),
	}
}

// CreateFolder 创建收藏夹
func (r *favoriteRepository) CreateFolder(folder *models.FavoriteFolder) error {
	return r.db.Create(folder).Error
}

// GetFolderByID 根据ID获取收藏夹
func (r *favoriteRepository) GetFolderByID(id uint) (*models.FavoriteFolder, error) {
	var folder models.FavoriteFolder
	err := r.db.Preload("User").First(&folder, id).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

// GetUserFolders 获取用户的收藏夹列表
func (r *favoriteRepository) GetUserFolders(userID uint) ([]*models.FavoriteFolder, error) {
	var folders []*models.FavoriteFolder
	err := r.db.Where("user_id = ?", userID).
		Order("is_default DESC, sort ASC, created_at ASC").
		Find(&folders).Error
	return folders, err
}

// UpdateFolder 更新收藏夹
func (r *favoriteRepository) UpdateFolder(folder *models.FavoriteFolder) error {
	return r.db.Save(folder).Error
}

// DeleteFolder 删除收藏夹
func (r *favoriteRepository) DeleteFolder(id uint) error {
	return r.db.Delete(&models.FavoriteFolder{}, id).Error
}

// GetDefaultFolder 获取用户的默认收藏夹
func (r *favoriteRepository) GetDefaultFolder(userID uint) (*models.FavoriteFolder, error) {
	var folder models.FavoriteFolder
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).
		First(&folder).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

// FolderExists 检查收藏夹名称是否已存在
func (r *favoriteRepository) FolderExists(userID uint, name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.FavoriteFolder{}).
		Where("user_id = ? AND name = ?", userID, name).
		Count(&count).Error
	return count > 0, err
}

// AddFavorite 添加文章到收藏
func (r *favoriteRepository) AddFavorite(favorite *models.ArticleFavorite) error {
	return r.db.Create(favorite).Error
}

// RemoveFavorite 移除文章收藏
func (r *favoriteRepository) RemoveFavorite(userID, articleID uint) error {
	return r.db.Where("user_id = ? AND article_id = ?", userID, articleID).
		Delete(&models.ArticleFavorite{}).Error
}

// GetFavorite 获取收藏记录
func (r *favoriteRepository) GetFavorite(userID, articleID uint) (*models.ArticleFavorite, error) {
	var favorite models.ArticleFavorite
	err := r.db.Where("user_id = ? AND article_id = ?", userID, articleID).
		Preload("Folder").
		First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

// IsFavorited 检查文章是否已被用户收藏
func (r *favoriteRepository) IsFavorited(userID, articleID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.ArticleFavorite{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count).Error
	return count > 0, err
}

// GetFolderFavorites 获取收藏夹中的文章
func (r *favoriteRepository) GetFolderFavorites(folderID uint, page, size int) ([]*models.ArticleFavorite, int64, error) {
	var favorites []*models.ArticleFavorite
	var total int64

	// 计算总数
	if err := r.db.Model(&models.ArticleFavorite{}).
		Where("folder_id = ?", folderID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	err := r.db.Where("folder_id = ?", folderID).
		Preload("Article").
		Preload("Article.User").
		Preload("Article.Category").
		Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&favorites).Error

	return favorites, total, err
}

// GetUserFavorites 获取用户的所有收藏文章
func (r *favoriteRepository) GetUserFavorites(userID uint, page, size int) ([]*models.ArticleFavorite, int64, error) {
	var favorites []*models.ArticleFavorite
	var total int64

	// 计算总数
	if err := r.db.Model(&models.ArticleFavorite{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	err := r.db.Where("user_id = ?", userID).
		Preload("Article").
		Preload("Article.User").
		Preload("Article.Category").
		Preload("Folder").
		Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&favorites).Error

	return favorites, total, err
}

// MoveFavorite 移动收藏到其他收藏夹
func (r *favoriteRepository) MoveFavorite(userID, articleID, newFolderID uint) error {
	return r.db.Model(&models.ArticleFavorite{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Update("folder_id", newFolderID).Error
}

// GetFavoriteCount 获取文章的收藏数量
func (r *favoriteRepository) GetFavoriteCount(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.ArticleFavorite{}).
		Where("article_id = ?", articleID).
		Count(&count).Error
	return count, err
}