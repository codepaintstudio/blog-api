package controllers

import (
	"strconv"

	"blog-api/internal/services"
	"blog-api/pkg/response"
	"blog-api/pkg/errors"

	"github.com/gin-gonic/gin"
)

// LikeController 点赞控制器
type LikeController struct {
	likeService services.LikeService
}

// NewLikeController 创建点赞控制器实例
func NewLikeController(likeService services.LikeService) *LikeController {
	return &LikeController{
		likeService: likeService,
	}
}

// ToggleArticleLike 切换文章点赞状态
// @Summary 切换文章点赞状态
// @Description 用户对文章进行点赞或取消点赞操作
// @Tags 点赞
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response{data=services.ArticleLikeResult}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /articles/{id}/like [post]
// @Security BearerAuth
func (c *LikeController) ToggleArticleLike(ctx *gin.Context) {
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

	// 执行点赞操作
	result, err := c.likeService.ToggleArticleLike(userID.(uint), uint(articleID))
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	if result.IsLiked {
		response.SuccessWithMessage(ctx, result, "点赞成功")
	} else {
		response.SuccessWithMessage(ctx, result, "取消点赞成功")
	}
}

// ToggleCommentLike 切换评论点赞状态
// @Summary 切换评论点赞状态
// @Description 用户对评论进行点赞或取消点赞操作
// @Tags 点赞
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response{data=services.CommentLikeResult}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /comments/{id}/like [post]
// @Security BearerAuth
func (c *LikeController) ToggleCommentLike(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.ErrorWithMessage(ctx, response.UNAUTHORIZED, "用户未登录")
		return
	}

	// 获取评论ID
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的评论ID")
		return
	}

	// 执行点赞操作
	result, err := c.likeService.ToggleCommentLike(userID.(uint), uint(commentID))
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	if result.IsLiked {
		response.SuccessWithMessage(ctx, result, "点赞成功")
	} else {
		response.SuccessWithMessage(ctx, result, "取消点赞成功")
	}
}

// GetUserArticleLikes 获取用户点赞的文章列表
// @Summary 获取用户点赞的文章列表
// @Description 获取当前用户点赞的所有文章列表，支持分页
// @Tags 点赞
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageData}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /likes/articles [get]
// @Security BearerAuth
func (c *LikeController) GetUserArticleLikes(ctx *gin.Context) {
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

	// 获取用户点赞的文章列表
	likes, total, err := c.likeService.GetUserArticleLikes(userID.(uint), page, size)
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, "获取点赞列表失败")
		return
	}

	response.SuccessPage(ctx, likes, total, page, size)
}

// GetUserCommentLikes 获取用户点赞的评论列表
// @Summary 获取用户点赞的评论列表
// @Description 获取当前用户点赞的所有评论列表，支持分页
// @Tags 点赞
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageData}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /likes/comments [get]
// @Security BearerAuth
func (c *LikeController) GetUserCommentLikes(ctx *gin.Context) {
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

	// 获取用户点赞的评论列表
	likes, total, err := c.likeService.GetUserCommentLikes(userID.(uint), page, size)
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, "获取点赞列表失败")
		return
	}

	response.SuccessPage(ctx, likes, total, page, size)
}

// GetArticleLikers 获取文章的点赞用户列表
// @Summary 获取文章的点赞用户列表
// @Description 获取指定文章的所有点赞用户列表，支持分页
// @Tags 点赞
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageData}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /articles/{id}/likers [get]
func (c *LikeController) GetArticleLikers(ctx *gin.Context) {
	// 获取文章ID
	articleIDStr := ctx.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的文章ID")
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

	// 获取文章的点赞用户列表
	users, total, err := c.likeService.GetArticleLikers(uint(articleID), page, size)
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, "获取点赞用户列表失败")
		return
	}

	response.SuccessPage(ctx, users, total, page, size)
}