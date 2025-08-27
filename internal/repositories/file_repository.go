package repositories

import (
	"blog-api/internal/models"
	"blog-api/pkg/database"
	"gorm.io/gorm"
)

// FileRepository 文件仓库接口
type FileRepository interface {
	Create(file *models.File) error
	GetByID(id uint) (*models.File, error)
	GetByHash(hash string) (*models.File, error)
	GetUserFiles(userID uint, page, size int) ([]*models.File, int64, error)
	Update(file *models.File) error
	Delete(id uint) error
	UpdateRefCount(id uint, count int) error
	GetByPath(path string) (*models.File, error)
}

// fileRepository 文件仓库实现
type fileRepository struct {
	db *gorm.DB
}

// NewFileRepository 创建文件仓库实例
func NewFileRepository() FileRepository {
	return &fileRepository{
		db: database.GetDB(),
	}
}

// Create 创建文件记录
func (r *fileRepository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

// GetByID 根据ID获取文件
func (r *fileRepository) GetByID(id uint) (*models.File, error) {
	var file models.File
	err := r.db.Where("id = ?", id).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// GetByHash 根据哈希值获取文件
func (r *fileRepository) GetByHash(hash string) (*models.File, error) {
	var file models.File
	err := r.db.Where("hash = ?", hash).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// GetUserFiles 获取用户文件列表（分页）
func (r *fileRepository) GetUserFiles(userID uint, page, size int) ([]*models.File, int64, error) {
	var files []*models.File
	var total int64

	// 计算总数
	if err := r.db.Model(&models.File{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&files).Error

	return files, total, err
}

// Update 更新文件信息
func (r *fileRepository) Update(file *models.File) error {
	return r.db.Save(file).Error
}

// Delete 删除文件记录
func (r *fileRepository) Delete(id uint) error {
	return r.db.Delete(&models.File{}, id).Error
}

// UpdateRefCount 更新文件引用计数
func (r *fileRepository) UpdateRefCount(id uint, count int) error {
	return r.db.Model(&models.File{}).Where("id = ?", id).Update("ref_count", count).Error
}


// GetByPath 根据路径获取文件
func (r *fileRepository) GetByPath(path string) (*models.File, error) {
	var file models.File
	err := r.db.Where("path = ?", path).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}