package repositories

import (
	"blog-api/internal/models"
	"blog-api/pkg/database"
	
	"gorm.io/gorm"
)

// UserFilters 用户筛选条件
type UserFilters struct {
	Role         string `form:"role"`
	Status       string `form:"status"`
	EmailVerified *bool `form:"email_verified"`
	Keyword      string `form:"keyword"` // 用户名或邮箱关键词搜索
}

// UserRepository 用户仓库接口
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(offset, limit int) ([]*models.User, int64, error)
	GetUsersWithFilters(page, size int, filters *UserFilters) ([]*models.User, int64, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
	UpdateLoginTime(id uint) error
}

// userRepository 用户仓库实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建新的用户仓库
func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.GetDB(),
	}
}

// Create 创建用户
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户（软删除）
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// List 获取用户列表
func (r *userRepository) List(offset, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64
	
	// 获取总数
	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
}

// ExistsByEmail 检查邮箱是否已存在
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// ExistsByUsername 检查用户名是否已存在
func (r *userRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// UpdateLoginTime 更新用户登录时间
func (r *userRepository) UpdateLoginTime(id uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("login_at", gorm.Expr("NOW()")).Error
}

// GetUsersWithFilters 根据筛选条件获取用户列表
func (r *userRepository) GetUsersWithFilters(page, size int, filters *UserFilters) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64
	
	// 构建查询
	query := r.db.Model(&models.User{})
	
	// 应用筛选条件
	if filters != nil {
		if filters.Role != "" {
			query = query.Where("role = ?", filters.Role)
		}
		if filters.Status != "" {
			query = query.Where("status = ?", filters.Status)
		}
		if filters.EmailVerified != nil {
			query = query.Where("email_verified = ?", *filters.EmailVerified)
		}
		if filters.Keyword != "" {
			query = query.Where("username LIKE ? OR email LIKE ?", 
				"%"+filters.Keyword+"%", "%"+filters.Keyword+"%")
		}
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	offset := (page - 1) * size
	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
}