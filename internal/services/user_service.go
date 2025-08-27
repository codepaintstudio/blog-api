package services

import (
	"errors"
	"time"

	"blog-api/internal/models"
	"blog-api/internal/repositories"
	"blog-api/internal/utils"
	"blog-api/pkg/logger"
	
	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	Register(req *RegisterRequest) (*models.User, error)
	Login(req *LoginRequest) (*LoginResponse, error)
	GetProfile(userID uint) (*models.User, error)
	UpdateProfile(userID uint, req *UpdateProfileRequest) (*models.User, error)
	ChangePassword(userID uint, req *ChangePasswordRequest) error
	RefreshToken(refreshToken string) (*TokenResponse, error)
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Nickname string `json:"nickname" validate:"max=50"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Account  string `json:"account" validate:"required"`  // 用户名或邮箱
	Password string `json:"password" validate:"required"`
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" validate:"max=50"`
	Bio      string `json:"bio" validate:"max=500"`
	Avatar   string `json:"avatar"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	User         *UserResponse  `json:"user"`
	TokenResponse
}

// TokenResponse 令牌响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID            uint       `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	Nickname      string     `json:"nickname"`
	Avatar        string     `json:"avatar"`
	Bio           string     `json:"bio"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LoginAt       *time.Time `json:"login_at"`
}

// userService 用户服务实现
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService 创建新的用户服务
func NewUserService() UserService {
	return &userService{
		userRepo: repositories.NewUserRepository(),
	}
}

// Register 用户注册
func (s *userService) Register(req *RegisterRequest) (*models.User, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		logger.Error("检查用户名失败", logger.Err("error", err))
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}
	
	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		logger.Error("检查邮箱失败", logger.Err("error", err))
		return nil, err
	}
	if exists {
		return nil, errors.New("邮箱已被注册")
	}
	
	// 验证密码格式
	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, err
	}
	
	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Error("密码加密失败", logger.Err("error", err))
		return nil, err
	}
	
	// 创建用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Role:     "user",
		Status:   "active",
	}
	
	if err := s.userRepo.Create(user); err != nil {
		logger.Error("创建用户失败", logger.Err("error", err))
		return nil, err
	}
	
	logger.Info("用户注册成功",
		logger.String("username", user.Username),
		logger.String("email", user.Email),
		logger.Int("user_id", int(user.ID)),
	)
	
	return user, nil
}

// Login 用户登录
func (s *userService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 根据账户查找用户（支持用户名或邮箱登录）
	var user *models.User
	var err error
	
	// 尝试按邮箱查找
	if utils.IsEmail(req.Account) {
		user, err = s.userRepo.GetByEmail(req.Account)
	} else {
		user, err = s.userRepo.GetByUsername(req.Account)
	}
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		logger.Error("查找用户失败", logger.Err("error", err))
		return nil, err
	}
	
	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.New("用户账号已被禁用")
	}
	
	// 验证密码
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		logger.Warn("用户登录密码错误",
			logger.String("account", req.Account),
			logger.Int("user_id", int(user.ID)),
		)
		return nil, errors.New("密码错误")
	}
	
	// 生成令牌
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		logger.Error("生成令牌失败", logger.Err("error", err))
		return nil, err
	}
	
	// 更新登录时间
	if err := s.userRepo.UpdateLoginTime(user.ID); err != nil {
		logger.Warn("更新登录时间失败", logger.Err("error", err))
	}
	
	logger.Info("用户登录成功",
		logger.String("username", user.Username),
		logger.Int("user_id", int(user.ID)),
	)
	
	return &LoginResponse{
		User: s.toUserResponse(user),
		TokenResponse: TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    3600, // 1小时，应该从配置读取
		},
	}, nil
}

// GetProfile 获取用户资料
func (s *userService) GetProfile(userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return user, nil
}

// UpdateProfile 更新用户资料
func (s *userService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	
	if err := s.userRepo.Update(user); err != nil {
		logger.Error("更新用户资料失败", logger.Err("error", err))
		return nil, err
	}
	
	logger.Info("用户资料更新成功",
		logger.Int("user_id", int(userID)),
		logger.String("username", user.Username),
	)
	
	return user, nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}
	
	// 验证旧密码
	if err := utils.CheckPassword(user.Password, req.OldPassword); err != nil {
		return errors.New("原密码错误")
	}
	
	// 验证新密码格式
	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		return err
	}
	
	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		logger.Error("新密码加密失败", logger.Err("error", err))
		return err
	}
	
	// 更新密码
	user.Password = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		logger.Error("更新密码失败", logger.Err("error", err))
		return err
	}
	
	logger.Info("用户密码修改成功",
		logger.Int("user_id", int(userID)),
		logger.String("username", user.Username),
	)
	
	return nil
}

// RefreshToken 刷新令牌
func (s *userService) RefreshToken(refreshToken string) (*TokenResponse, error) {
	// 验证刷新令牌
	if err := utils.ValidateTokenType(refreshToken, "refresh"); err != nil {
		return nil, err
	}
	
	// 刷新访问令牌
	newAccessToken, err := utils.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	
	return &TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken, // 刷新令牌保持不变
		ExpiresIn:    3600,         // 1小时，应该从配置读取
	}, nil
}

// toUserResponse 转换为用户响应格式
func (s *userService) toUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Nickname:      user.Nickname,
		Avatar:        user.Avatar,
		Bio:           user.Bio,
		Role:          user.Role,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		LoginAt:       user.LoginAt,
	}
}