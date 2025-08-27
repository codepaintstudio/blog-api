package controllers

import (
	"strconv"

	"blog-api/internal/repositories"
	"blog-api/internal/services"
	"blog-api/pkg/response"
	"blog-api/pkg/errors"

	"github.com/gin-gonic/gin"
)

// AdminController 管理员控制器
type AdminController struct {
	adminService services.AdminService
}

// NewAdminController 创建管理员控制器实例
func NewAdminController(adminService services.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
	}
}

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 管理员获取用户列表，支持筛选和分页
// @Tags 管理员
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Param role query string false "用户角色" Enums(admin,user)
// @Param status query string false "用户状态" Enums(active,inactive)
// @Param email_verified query bool false "邮箱验证状态"
// @Param keyword query string false "搜索关键词（用户名或邮箱）"
// @Success 200 {object} response.Response{data=response.PageData}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/users [get]
// @Security BearerAuth
func (c *AdminController) GetUserList(ctx *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	// 获取筛选条件
	var filters repositories.UserFilters
	if err := ctx.ShouldBindQuery(&filters); err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "参数错误")
		return
	}

	// 获取用户列表
	users, total, err := c.adminService.GetUserList(page, size, &filters)
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.SuccessPage(ctx, users, total, page, size)
}

// GetUserDetails 获取用户详情
// @Summary 获取用户详情
// @Description 管理员获取指定用户的详细信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=services.AdminUserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/users/{id} [get]
// @Security BearerAuth
func (c *AdminController) GetUserDetails(ctx *gin.Context) {
	// 获取用户ID
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的用户ID")
		return
	}

	// 获取用户详情
	user, err := c.adminService.GetUserDetails(uint(userID))
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	response.Success(ctx, user)
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Description 管理员更新用户状态（激活/禁用）
// @Tags 管理员
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param status body object true "状态信息" example({"status":"active"})
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/users/{id}/status [put]
// @Security BearerAuth
func (c *AdminController) UpdateUserStatus(ctx *gin.Context) {
	// 获取用户ID
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的用户ID")
		return
	}

	// 绑定请求参数
	var req struct {
		Status string `json:"status" validate:"required,oneof=active inactive"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "参数错误")
		return
	}

	// 更新用户状态
	err = c.adminService.UpdateUserStatus(uint(userID), req.Status)
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, nil, "用户状态更新成功")
}

// UpdateUserRole 更新用户角色
// @Summary 更新用户角色
// @Description 管理员更新用户角色
// @Tags 管理员
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param role body object true "角色信息" example({"role":"admin"})
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/users/{id}/role [put]
// @Security BearerAuth
func (c *AdminController) UpdateUserRole(ctx *gin.Context) {
	// 获取用户ID
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的用户ID")
		return
	}

	// 绑定请求参数
	var req struct {
		Role string `json:"role" validate:"required,oneof=admin user"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "参数错误")
		return
	}

	// 更新用户角色
	err = c.adminService.UpdateUserRole(uint(userID), req.Role)
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, nil, "用户角色更新成功")
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 管理员删除用户（软删除）
// @Tags 管理员
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/users/{id} [delete]
// @Security BearerAuth
func (c *AdminController) DeleteUser(ctx *gin.Context) {
	// 获取用户ID
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的用户ID")
		return
	}

	// 删除用户
	err = c.adminService.DeleteUser(uint(userID))
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, nil, "用户删除成功")
}

// GetSystemStats 获取系统统计
// @Summary 获取系统统计
// @Description 管理员获取系统整体统计信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=services.SystemStats}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/stats/system [get]
// @Security BearerAuth
func (c *AdminController) GetSystemStats(ctx *gin.Context) {
	stats, err := c.adminService.GetSystemStats()
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(ctx, stats)
}

// GetUserStats 获取用户统计
// @Summary 获取用户统计
// @Description 管理员获取用户相关统计信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Param days query int false "统计天数" default(30)
// @Success 200 {object} response.Response{data=services.UserStats}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/stats/users [get]
// @Security BearerAuth
func (c *AdminController) GetUserStats(ctx *gin.Context) {
	days, _ := strconv.Atoi(ctx.DefaultQuery("days", "30"))
	if days <= 0 || days > 365 {
		days = 30
	}

	stats, err := c.adminService.GetUserStats(days)
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(ctx, stats)
}

// GetContentStats 获取内容统计
// @Summary 获取内容统计
// @Description 管理员获取内容相关统计信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Param days query int false "统计天数" default(30)
// @Success 200 {object} response.Response{data=services.ContentStats}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/stats/content [get]
// @Security BearerAuth
func (c *AdminController) GetContentStats(ctx *gin.Context) {
	days, _ := strconv.Atoi(ctx.DefaultQuery("days", "30"))
	if days <= 0 || days > 365 {
		days = 30
	}

	stats, err := c.adminService.GetContentStats(days)
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(ctx, stats)
}