package controllers

import (
	"blog-api/internal/middlewares"
	"blog-api/internal/services"
	"blog-api/pkg/logger"
	"blog-api/pkg/response"
	
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserController 用户控制器
type UserController struct {
	userService services.UserService
	validator   *validator.Validate
}

// NewUserController 创建新的用户控制器
func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
		validator:   validator.New(),
	}
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前登录用户的详细资料
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=services.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/profile [get]
func (uc *UserController) GetProfile(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := middlewares.GetCurrentUserID(c)
	if !exists {
		response.Error(c, response.UNAUTHORIZED)
		return
	}
	
	// 调用服务层获取用户资料
	user, err := uc.userService.GetProfile(userID)
	if err != nil {
		logger.Error("获取用户资料失败",
			logger.Err("error", err),
			logger.Int("user_id", int(userID)),
		)
		
		if err.Error() == "用户不存在" {
			response.Error(c, response.NOT_FOUND)
		} else {
			response.ErrorWithMessage(c, response.ERROR, "获取用户资料失败")
		}
		return
	}
	
	// 转换为响应格式
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
	
	response.Success(c, userResponse)
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新当前登录用户的资料信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body services.UpdateProfileRequest true "更新信息"
// @Success 200 {object} response.Response{data=services.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /users/profile [put]
func (uc *UserController) UpdateProfile(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := middlewares.GetCurrentUserID(c)
	if !exists {
		response.Error(c, response.UNAUTHORIZED)
		return
	}
	
	var req services.UpdateProfileRequest
	
	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("更新用户资料参数绑定失败",
			logger.Err("error", err),
			logger.Int("user_id", int(userID)),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "参数格式错误")
		return
	}
	
	// 验证参数
	if err := uc.validator.Struct(&req); err != nil {
		logger.Warn("更新用户资料参数验证失败",
			logger.Err("error", err),
			logger.Int("user_id", int(userID)),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "参数验证失败")
		return
	}
	
	// 调用服务层更新用户资料
	user, err := uc.userService.UpdateProfile(userID, &req)
	if err != nil {
		logger.Error("更新用户资料失败",
			logger.Err("error", err),
			logger.Int("user_id", int(userID)),
		)
		
		if err.Error() == "用户不存在" {
			response.Error(c, response.NOT_FOUND)
		} else {
			response.ErrorWithMessage(c, response.ERROR, "更新用户资料失败")
		}
		return
	}
	
	// 转换为响应格式
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
	
	response.SuccessWithMessage(c, userResponse, "资料更新成功")
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body services.ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /users/password [put]
func (uc *UserController) ChangePassword(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := middlewares.GetCurrentUserID(c)
	if !exists {
		response.Error(c, response.UNAUTHORIZED)
		return
	}
	
	var req services.ChangePasswordRequest
	
	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("修改密码参数绑定失败",
			logger.Err("error", err),
			logger.Int("user_id", int(userID)),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "参数格式错误")
		return
	}
	
	// 验证参数
	if err := uc.validator.Struct(&req); err != nil {
		logger.Warn("修改密码参数验证失败",
			logger.Err("error", err),
			logger.Int("user_id", int(userID)),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "参数验证失败")
		return
	}
	
	// 调用服务层修改密码
	if err := uc.userService.ChangePassword(userID, &req); err != nil {
		logger.Error("修改密码失败",
			logger.Err("error", err),
			logger.Int("user_id", int(userID)),
		)
		
		if err.Error() == "用户不存在" {
			response.Error(c, response.NOT_FOUND)
		} else if err.Error() == "原密码错误" {
			response.ErrorWithMessage(c, response.UNAUTHORIZED, err.Error())
		} else {
			response.ErrorWithMessage(c, response.ERROR, "修改密码失败")
		}
		return
	}
	
	response.SuccessWithMessage(c, nil, "密码修改成功")
}