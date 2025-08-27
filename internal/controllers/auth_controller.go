package controllers

import (
	"blog-api/internal/services"
	"blog-api/pkg/logger"
	"blog-api/pkg/response"
	
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthController 认证控制器
type AuthController struct {
	userService services.UserService
	validator   *validator.Validate
}

// NewAuthController 创建新的认证控制器
func NewAuthController(userService services.UserService) *AuthController {
	return &AuthController{
		userService: userService,
		validator:   validator.New(),
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口，创建新用户账户
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body services.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=services.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req services.RegisterRequest
	
	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("用户注册参数绑定失败",
			logger.Err("error", err),
			logger.String("ip", c.ClientIP()),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "参数格式错误")
		return
	}
	
	// 验证参数
	if err := ac.validator.Struct(&req); err != nil {
		logger.Warn("用户注册参数验证失败",
			logger.Err("error", err),
			logger.String("username", req.Username),
			logger.String("email", req.Email),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "参数验证失败")
		return
	}
	
	// 调用服务层注册用户
	user, err := ac.userService.Register(&req)
	if err != nil {
		logger.Error("用户注册失败",
			logger.Err("error", err),
			logger.String("username", req.Username),
			logger.String("email", req.Email),
		)
		
		// 根据错误类型返回不同的响应
		if err.Error() == "用户名已存在" || err.Error() == "邮箱已被注册" {
			response.ErrorWithMessage(c, response.ALREADY_EXISTS, err.Error())
		} else {
			response.ErrorWithMessage(c, response.ERROR, "注册失败，请稍后重试")
		}
		return
	}
	
	// 返回成功响应（不包含密码）
	userResponse := &services.UserResponse{
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
	
	response.SuccessWithMessage(c, userResponse, "注册成功")
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，支持用户名或邮箱登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body services.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=services.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req services.LoginRequest
	
	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("用户登录参数绑定失败",
			logger.Err("error", err),
			logger.String("ip", c.ClientIP()),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "参数格式错误")
		return
	}
	
	// 验证参数
	if err := ac.validator.Struct(&req); err != nil {
		logger.Warn("用户登录参数验证失败",
			logger.Err("error", err),
			logger.String("account", req.Account),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "参数验证失败")
		return
	}
	
	// 调用服务层登录
	loginResp, err := ac.userService.Login(&req)
	if err != nil {
		logger.Warn("用户登录失败",
			logger.Err("error", err),
			logger.String("account", req.Account),
			logger.String("ip", c.ClientIP()),
		)
		
		// 根据错误类型返回不同的响应
		if err.Error() == "用户不存在" || err.Error() == "密码错误" {
			response.ErrorWithMessage(c, response.UNAUTHORIZED, "用户名或密码错误")
		} else if err.Error() == "用户账号已被禁用" {
			response.ErrorWithMessage(c, response.FORBIDDEN, err.Error())
		} else {
			response.ErrorWithMessage(c, response.ERROR, "登录失败，请稍后重试")
		}
		return
	}
	
	response.SuccessWithMessage(c, loginResp, "登录成功")
}

// RefreshToken 刷新令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body map[string]string true "刷新令牌" example({"refresh_token": "your_refresh_token"})
// @Success 200 {object} response.Response{data=services.TokenResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/refresh [post]
func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	
	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "参数格式错误")
		return
	}
	
	// 验证参数
	if err := ac.validator.Struct(&req); err != nil {
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "刷新令牌不能为空")
		return
	}
	
	// 调用服务层刷新令牌
	tokenResp, err := ac.userService.RefreshToken(req.RefreshToken)
	if err != nil {
		logger.Warn("刷新令牌失败",
			logger.Err("error", err),
			logger.String("ip", c.ClientIP()),
		)
		response.ErrorWithMessage(c, response.UNAUTHORIZED, "刷新令牌无效或已过期")
		return
	}
	
	response.SuccessWithMessage(c, tokenResp, "令牌刷新成功")
}