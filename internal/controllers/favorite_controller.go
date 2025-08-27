package controllers

import (
	"strconv"

	"blog-api/internal/services"
	"blog-api/pkg/errors"
	"blog-api/pkg/response"

	"github.com/gin-gonic/gin"
)

// FavoriteController 收藏控制器
type FavoriteController struct {
	favoriteService services.FavoriteService
}

// NewFavoriteController 创建收藏控制器实例
func NewFavoriteController(favoriteService services.FavoriteService) *FavoriteController {
	return &FavoriteController{
		favoriteService: favoriteService,
	}
}

// CreateFolder 创建收藏夹
// @Summary 创建收藏夹
// @Description 用户创建新的收藏夹
// @Tags 收藏
// @Accept json
// @Produce json
// @Param folder body services.CreateFolderRequest true "收藏夹信息"
// @Success 200 {object} response.Response{data=models.FavoriteFolder}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /favorites/folders [post]
// @Security BearerAuth
func (c *FavoriteController) CreateFolder(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.ErrorWithMessage(ctx, response.UNAUTHORIZED, "用户未登录")
		return
	}

	// 绑定请求参数
	var req services.CreateFolderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "参数错误")
		return
	}

	// 创建收藏夹
	folder, err := c.favoriteService.CreateFolder(userID.(uint), &req)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	response.SuccessWithMessage(ctx, folder, "收藏夹创建成功")
}

// UpdateFolder 更新收藏夹
// @Summary 更新收藏夹
// @Description 用户更新收藏夹信息
// @Tags 收藏
// @Accept json
// @Produce json
// @Param id path int true "收藏夹ID"
// @Param folder body services.UpdateFolderRequest true "收藏夹信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /favorites/folders/{id} [put]
// @Security BearerAuth
func (c *FavoriteController) UpdateFolder(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.ErrorWithMessage(ctx, response.UNAUTHORIZED, "用户未登录")
		return
	}

	// 获取收藏夹ID
	folderIDStr := ctx.Param("id")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的收藏夹ID")
		return
	}

	// 绑定请求参数
	var req services.UpdateFolderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "参数错误")
		return
	}

	// 更新收藏夹
	_, err = c.favoriteService.UpdateFolder(userID.(uint), uint(folderID), &req)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	response.SuccessWithMessage(ctx, nil, "收藏夹更新成功")
}

// DeleteFolder 删除收藏夹
// @Summary 删除收藏夹
// @Description 用户删除收藏夹（会同时删除其中的收藏记录）
// @Tags 收藏
// @Accept json
// @Produce json
// @Param id path int true "收藏夹ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /favorites/folders/{id} [delete]
// @Security BearerAuth
func (c *FavoriteController) DeleteFolder(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.ErrorWithMessage(ctx, response.UNAUTHORIZED, "用户未登录")
		return
	}

	// 获取收藏夹ID
	folderIDStr := ctx.Param("id")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的收藏夹ID")
		return
	}

	// 删除收藏夹
	err = c.favoriteService.DeleteFolder(userID.(uint), uint(folderID))
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	response.SuccessWithMessage(ctx, nil, "收藏夹删除成功")
}

// GetUserFolders 获取用户收藏夹列表
// @Summary 获取用户收藏夹列表
// @Description 获取当前用户的所有收藏夹列表
// @Tags 收藏
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]models.FavoriteFolder}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /favorites/folders [get]
// @Security BearerAuth
func (c *FavoriteController) GetUserFolders(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.ErrorWithMessage(ctx, response.UNAUTHORIZED, "用户未登录")
		return
	}

	// 获取用户收藏夹列表
	folders, err := c.favoriteService.GetUserFolders(userID.(uint))
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, "获取收藏夹列表失败")
		return
	}

	response.Success(ctx, folders)
}

// AddToFavorite 添加文章到收藏夹
// @Summary 添加文章到收藏夹
// @Description 将文章添加到指定收藏夹
// @Tags 收藏
// @Accept json
// @Produce json
// @Param favorite body services.FavoriteAddRequest true "收藏信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /favorites [post]
// @Security BearerAuth
func (c *FavoriteController) AddToFavorite(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.ErrorWithMessage(ctx, response.UNAUTHORIZED, "用户未登录")
		return
	}

	// 绑定请求参数
	var req services.FavoriteAddRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "参数错误")
		return
	}

	// 添加到收藏夹
	err := c.favoriteService.AddToFavorite(userID.(uint), &req)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	response.SuccessWithMessage(ctx, nil, "收藏成功")
}

// RemoveFromFavorite 从收藏夹移除文章
// @Summary 从收藏夹移除文章
// @Description 将文章从指定收藏夹中移除
// @Tags 收藏
// @Accept json
// @Produce json
// @Param article_id path int true "文章ID"
// @Param folder_id query int false "收藏夹ID（不指定则从所有收藏夹移除）"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /favorites/{article_id} [delete]
// @Security BearerAuth
func (c *FavoriteController) RemoveFromFavorite(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.ErrorWithMessage(ctx, response.UNAUTHORIZED, "用户未登录")
		return
	}

	// 获取文章ID
	articleIDStr := ctx.Param("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的文章ID")
		return
	}

	// 获取收藏夹ID（可选）
	var folderID *uint
	if folderIDStr := ctx.Query("folder_id"); folderIDStr != "" {
		if id, err := strconv.ParseUint(folderIDStr, 10, 32); err == nil {
			temp := uint(id)
			folderID = &temp
		}
	}

	// 从收藏夹移除
	err = c.favoriteService.RemoveFromFavorite(userID.(uint), uint(articleID), folderID)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	response.SuccessWithMessage(ctx, nil, "取消收藏成功")
}

// GetFolderArticles 获取收藏夹中的文章列表
// @Summary 获取收藏夹中的文章列表
// @Description 获取指定收藏夹中的所有文章，支持分页
// @Tags 收藏
// @Accept json
// @Produce json
// @Param id path int true "收藏夹ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageData}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /favorites/folders/{id}/articles [get]
// @Security BearerAuth
func (c *FavoriteController) GetFolderArticles(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.ErrorWithMessage(ctx, response.UNAUTHORIZED, "用户未登录")
		return
	}

	// 获取收藏夹ID
	folderIDStr := ctx.Param("id")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的收藏夹ID")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	// 获取收藏夹文章列表
	articles, total, err := c.favoriteService.GetFolderArticles(userID.(uint), uint(folderID), page, size)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	response.SuccessPage(ctx, articles, total, page, size)
}

// GetUserFavorites 获取用户所有收藏的文章
// @Summary 获取用户所有收藏的文章
// @Description 获取当前用户所有收藏的文章列表，支持分页
// @Tags 收藏
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageData}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /favorites/articles [get]
// @Security BearerAuth
func (c *FavoriteController) GetUserFavorites(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.ErrorWithMessage(ctx, response.UNAUTHORIZED, "用户未登录")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	// 获取用户收藏的文章列表
	articles, total, err := c.favoriteService.GetUserFavorites(userID.(uint), page, size)
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, "获取收藏列表失败")
		return
	}

	response.SuccessPage(ctx, articles, total, page, size)
}

// CheckFavoriteStatus 检查文章收藏状态
// @Summary 检查文章收藏状态
// @Description 检查指定文章是否已被当前用户收藏
// @Tags 收藏
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response{data=services.FavoriteStatusResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /favorites/{id}/status [get]
// @Security BearerAuth
func (c *FavoriteController) CheckFavoriteStatus(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.ErrorWithMessage(ctx, response.UNAUTHORIZED, "用户未登录")
		return
	}

	// 获取文章ID
	articleIDStr := ctx.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的文章ID")
		return
	}

	// 检查收藏状态
	status, err := c.favoriteService.CheckFavoriteStatus(userID.(uint), uint(articleID))
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, "检查收藏状态失败")
		return
	}

	response.Success(ctx, status)
}
