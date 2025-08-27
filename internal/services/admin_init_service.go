package services

import (
	"blog-api/internal/models"
	"blog-api/internal/repositories"
	"blog-api/internal/utils"
	"blog-api/pkg/config"
	"blog-api/pkg/logger"
)

// AdminInitService 管理员初始化服务
type AdminInitService struct {
	userRepo repositories.UserRepository
}

// NewAdminInitService 创建管理员初始化服务
func NewAdminInitService(userRepo repositories.UserRepository) *AdminInitService {
	return &AdminInitService{
		userRepo: userRepo,
	}
}

// InitializeAdmin 初始化系统管理员
func (s *AdminInitService) InitializeAdmin() error {
	// 检查是否已有管理员
	hasAdmin, err := s.userRepo.HasAdmin()
	if err != nil {
		return err
	}

	if hasAdmin {
		logger.Info("系统管理员已存在，跳过初始化")
		return nil
	}

	logger.Info("检测到系统中无管理员，开始创建默认管理员...")

	// 获取配置中的管理员信息
	cfg := config.GetConfig()
	if cfg == nil || cfg.Admin.Username == "" {
		logger.Error("管理员配置信息不完整")
		return nil
	}

	// 检查用户名和邮箱是否已被占用
	existsByUsername, err := s.userRepo.ExistsByUsername(cfg.Admin.Username)
	if err != nil {
		return err
	}
	if existsByUsername {
		logger.Warn("管理员用户名已被占用", logger.String("username", cfg.Admin.Username))
		return nil
	}

	existsByEmail, err := s.userRepo.ExistsByEmail(cfg.Admin.Email)
	if err != nil {
		return err
	}
	if existsByEmail {
		logger.Warn("管理员邮箱已被占用", logger.String("email", cfg.Admin.Email))
		return nil
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(cfg.Admin.Password)
	if err != nil {
		logger.Error("管理员密码加密失败", logger.Err("error", err))
		return err
	}

	// 创建管理员用户
	admin := &models.User{
		Username:      cfg.Admin.Username,
		Email:         cfg.Admin.Email,
		Nickname:      cfg.Admin.Nickname,
		Password:      hashedPassword,
		Role:          "admin",
		Status:        "active",
		EmailVerified: true,
	}

	err = s.userRepo.Create(admin)
	if err != nil {
		logger.Error("创建管理员失败", logger.Err("error", err))
		return err
	}

	logger.Info("系统管理员创建成功",
		logger.String("username", admin.Username),
		logger.String("email", admin.Email),
		logger.String("nickname", admin.Nickname),
	)

	return nil
}
