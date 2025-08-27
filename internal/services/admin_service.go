package services

import (
	"errors"
	"time"

	"blog-api/internal/models"
	"blog-api/internal/repositories"
	"blog-api/pkg/logger"

	"gorm.io/gorm"
)

// AdminService 管理员服务接口
type AdminService interface {
	// 用户管理
	GetUserList(page, size int, filters *repositories.UserFilters) ([]*AdminUserResponse, int64, error)
	GetUserDetails(userID uint) (*AdminUserResponse, error)
	UpdateUserStatus(userID uint, status string) error
	UpdateUserRole(userID uint, role string) error
	DeleteUser(userID uint) error
	
	// 系统统计
	GetSystemStats() (*SystemStats, error)
	GetUserStats(days int) (*UserStats, error)
	GetContentStats(days int) (*ContentStats, error)
}

// AdminUserResponse 管理员用户响应
type AdminUserResponse struct {
	ID            uint       `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	Nickname      string     `json:"nickname"`
	Avatar        string     `json:"avatar"`
	Bio           string     `json:"bio"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	LoginAt       *time.Time `json:"login_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	
	// 统计信息
	ArticleCount  int64 `json:"article_count"`
	CommentCount  int64 `json:"comment_count"`
	LikeCount     int64 `json:"like_count"`
	FavoriteCount int64 `json:"favorite_count"`
}

// SystemStats 系统统计
type SystemStats struct {
	TotalUsers       int64 `json:"total_users"`
	TotalArticles    int64 `json:"total_articles"`
	TotalComments    int64 `json:"total_comments"`
	TotalLikes       int64 `json:"total_likes"`
	TotalFavorites   int64 `json:"total_favorites"`
	TotalFiles       int64 `json:"total_files"`
	ActiveUsers      int64 `json:"active_users"`      // 最近30天活跃用户
	PublishedArticles int64 `json:"published_articles"`
	PendingComments  int64 `json:"pending_comments"`
}

// UserStats 用户统计
type UserStats struct {
	NewUsers       []DailyCount `json:"new_users"`
	ActiveUsers    []DailyCount `json:"active_users"`
	TotalUsers     int64        `json:"total_users"`
	GrowthRate     float64      `json:"growth_rate"`
}

// ContentStats 内容统计
type ContentStats struct {
	NewArticles     []DailyCount `json:"new_articles"`
	NewComments     []DailyCount `json:"new_comments"`
	TotalArticles   int64        `json:"total_articles"`
	TotalComments   int64        `json:"total_comments"`
	ArticleGrowthRate float64    `json:"article_growth_rate"`
	CommentGrowthRate float64    `json:"comment_growth_rate"`
}

// DailyCount 每日统计
type DailyCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// adminService 管理员服务实现
type adminService struct {
	userRepo     repositories.UserRepository
	articleRepo  repositories.ArticleRepository
	commentRepo  repositories.CommentRepository
	favoriteRepo repositories.FavoriteRepository
	fileRepo     repositories.FileRepository
}

// NewAdminService 创建管理员服务实例
func NewAdminService() AdminService {
	return &adminService{
		userRepo:     repositories.NewUserRepository(),
		articleRepo:  repositories.NewArticleRepository(),
		commentRepo:  repositories.NewCommentRepository(),
		favoriteRepo: repositories.NewFavoriteRepository(),
		fileRepo:     repositories.NewFileRepository(),
	}
}

// GetUserList 获取用户列表
func (s *adminService) GetUserList(page, size int, filters *repositories.UserFilters) ([]*AdminUserResponse, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	users, total, err := s.userRepo.GetUsersWithFilters(page, size, filters)
	if err != nil {
		logger.Error("获取用户列表失败", logger.Err("error", err))
		return nil, 0, errors.New("获取用户列表失败")
	}

	var responses []*AdminUserResponse
	for _, user := range users {
		response := s.convertToAdminUserResponse(user)
		responses = append(responses, response)
	}

	return responses, total, nil
}

// GetUserDetails 获取用户详情
func (s *adminService) GetUserDetails(userID uint) (*AdminUserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		logger.Error("获取用户详情失败", logger.Err("error", err))
		return nil, errors.New("获取用户详情失败")
	}

	response := s.convertToAdminUserResponse(user)
	return response, nil
}

// UpdateUserStatus 更新用户状态
func (s *adminService) UpdateUserStatus(userID uint, status string) error {
	if status != "active" && status != "inactive" {
		return errors.New("无效的用户状态")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return errors.New("获取用户信息失败")
	}

	user.Status = status
	if err := s.userRepo.Update(user); err != nil {
		logger.Error("更新用户状态失败", logger.Err("error", err))
		return errors.New("更新用户状态失败")
	}

	logger.Info("更新用户状态成功", 
		logger.Uint("user_id", userID),
		logger.String("status", status),
	)

	return nil
}

// UpdateUserRole 更新用户角色
func (s *adminService) UpdateUserRole(userID uint, role string) error {
	if role != "admin" && role != "user" {
		return errors.New("无效的用户角色")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return errors.New("获取用户信息失败")
	}

	user.Role = role
	if err := s.userRepo.Update(user); err != nil {
		logger.Error("更新用户角色失败", logger.Err("error", err))
		return errors.New("更新用户角色失败")
	}

	logger.Info("更新用户角色成功", 
		logger.Uint("user_id", userID),
		logger.String("role", role),
	)

	return nil
}

// DeleteUser 删除用户
func (s *adminService) DeleteUser(userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return errors.New("获取用户信息失败")
	}

	// 检查是否为管理员账户
	if user.Role == "admin" {
		return errors.New("不能删除管理员账户")
	}

	if err := s.userRepo.Delete(userID); err != nil {
		logger.Error("删除用户失败", logger.Err("error", err))
		return errors.New("删除用户失败")
	}

	logger.Info("删除用户成功", logger.Uint("user_id", userID))

	return nil
}

// GetSystemStats 获取系统统计
func (s *adminService) GetSystemStats() (*SystemStats, error) {
	stats := &SystemStats{}

	// 获取各项统计数据
	// 这里需要实现具体的统计逻辑，现在先返回基本结构
	// TODO: 实现具体的统计查询

	return stats, nil
}

// GetUserStats 获取用户统计
func (s *adminService) GetUserStats(days int) (*UserStats, error) {
	if days <= 0 || days > 365 {
		days = 30
	}

	stats := &UserStats{}
	
	// TODO: 实现用户统计逻辑

	return stats, nil
}

// GetContentStats 获取内容统计
func (s *adminService) GetContentStats(days int) (*ContentStats, error) {
	if days <= 0 || days > 365 {
		days = 30
	}

	stats := &ContentStats{}
	
	// TODO: 实现内容统计逻辑

	return stats, nil
}

// convertToAdminUserResponse 转换为管理员用户响应格式
func (s *adminService) convertToAdminUserResponse(user *models.User) *AdminUserResponse {
	response := &AdminUserResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Nickname:      user.Nickname,
		Avatar:        user.Avatar,
		Bio:           user.Bio,
		Role:          user.Role,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		LoginAt:       user.LoginAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}

	// TODO: 获取用户的统计信息
	// response.ArticleCount = ...
	// response.CommentCount = ...
	// response.LikeCount = ...
	// response.FavoriteCount = ...

	return response
}