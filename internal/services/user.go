package services

import (
	"errors"
	"fmt"
	"server/internal/models"
	"server/pkg/utils"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// 用户注册
func (s *UserService) Register(request *models.RegisterRequest) (*models.BaseUser, error) {
	// 验证密码是否匹配
	if request.Password != request.RepeatPassword {
		return nil, errors.New("passwords do not match")
	}

	// 检查邮箱是否已存在
	var existingUser models.User
	result := s.db.Where("email = ?", request.Email).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("email already exists")
	}

	// 检查用户名是否已存在
	result = s.db.Where("username = ?", request.Username).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("username already exists")
	}

	// 对密码进行哈希处理
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建新用户
	user := &models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
	}

	// 保存到数据库
	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 返回基础用户信息
	return &models.BaseUser{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// 用户登录
func (s *UserService) Login(request *models.LoginRequest) (*models.LoginResponse, error) {
	// 查找用户
	var user models.User
	result := s.db.Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		return nil, errors.New("invalid email or password")
	}

	// 验证密码
	if !utils.ValidatePassword(request.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 返回登录响应
	return &models.LoginResponse{
		Token: token,
		User: models.BaseUser{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

// GetUserById 根据用户ID获取用户信息
func (s *UserService) GetUserById(userId int) (*models.BaseUser, error) {
	var user models.User
	if err := s.db.First(&user, userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &models.BaseUser{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// GetUserDetail 获取用户详情（包含文章统计）
func (s *UserService) GetUserDetail(userId int) (*models.UserDetailResponse, error) {
	var user models.User
	if err := s.db.First(&user, userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// 统计用户文章数
	var articleCount int64
	if err := s.db.Model(&models.Article{}).Where("user_id = ?", userId).Count(&articleCount).Error; err != nil {
		return nil, err
	}

	return &models.UserDetailResponse{
		BaseUser: models.BaseUser{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
		},
		ArticleCount: int(articleCount),
	}, nil
}
