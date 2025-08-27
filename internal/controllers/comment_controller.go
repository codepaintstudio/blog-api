package controllers

import (
	"strconv"

	"blog-api/internal/services"
	"blog-api/pkg/response"
	"blog-api/pkg/errors"

	"github.com/gin-gonic/gin"
)

// CommentController 评论控制器
type CommentController struct {
	commentService services.CommentService
}

// NewCommentController 创建评论控制器实例
func NewCommentController(commentService services.CommentService) *CommentController {
	return &CommentController{
		commentService: commentService,
	}
}

// CreateComment 创建评论
// @Summary 创建评论
// @Description 用户对文章或其他评论进行评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param comment body services.CommentCreateRequest true "评论信息"
// @Success 200 {object} response.Response{data=models.Comment}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /comments [post]
// @Security BearerAuth
func (c *CommentController) CreateComment(ctx *gin.Context) {
	// 获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.ErrorWithMessage(ctx, response.UNAUTHORIZED, "用户未登录")
		return
	}

	// 绑定请求参数
	var req services.CommentCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "参数错误")
		return
	}

	// 获取客户端IP和User Agent
	clientIP := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	// 创建评论
	comment, err := c.commentService.CreateComment(userID.(uint), &req, clientIP, userAgent)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	response.SuccessWithMessage(ctx, comment, "评论发表成功")
}

// UpdateComment 更新评论
// @Summary 更新评论
// @Description 用户更新自己的评论内容
// @Tags 评论
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Param comment body services.CommentUpdateRequest true "评论信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /comments/{id} [put]
// @Security BearerAuth
func (c *CommentController) UpdateComment(ctx *gin.Context) {
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

	// 绑定请求参数
	var req services.CommentUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "参数错误")
		return
	}

	// 更新评论
	err = c.commentService.UpdateComment(userID.(uint), uint(commentID), &req)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	response.SuccessWithMessage(ctx, nil, "评论更新成功")
}

// DeleteComment 删除评论
// @Summary 删除评论
// @Description 用户删除自己的评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /comments/{id} [delete]
// @Security BearerAuth
func (c *CommentController) DeleteComment(ctx *gin.Context) {
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

	// 删除评论
	err = c.commentService.DeleteComment(userID.(uint), uint(commentID))
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	response.SuccessWithMessage(ctx, nil, "评论删除成功")
}

// GetComment 获取评论详情
// @Summary 获取评论详情
// @Description 获取指定评论的详细信息
// @Tags 评论
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response{data=services.CommentListResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /comments/{id} [get]
func (c *CommentController) GetComment(ctx *gin.Context) {
	// 获取评论ID
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的评论ID")
		return
	}

	// 获取当前用户ID（可选）
	var userID *uint
	if uid, exists := ctx.Get("user_id"); exists {
		temp := uid.(uint)
		userID = &temp
	}

	// 获取评论详情
	comment, err := c.commentService.GetCommentByID(uint(commentID), userID)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	response.Success(ctx, comment)
}

// GetArticleComments 获取文章评论列表
// @Summary 获取文章评论列表
// @Description 获取指定文章的评论列表，支持分页
// @Tags 评论
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageData}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /articles/{id}/comments [get]
func (c *CommentController) GetArticleComments(ctx *gin.Context) {
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

	// 获取当前用户ID（可选）
	var userID *uint
	if uid, exists := ctx.Get("user_id"); exists {
		temp := uid.(uint)
		userID = &temp
	}

	// 获取文章评论列表
	comments, total, err := c.commentService.GetArticleComments(uint(articleID), page, size, userID)
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, "获取评论列表失败")
		return
	}

	response.SuccessPage(ctx, comments, total, page, size)
}

// GetCommentReplies 获取评论回复列表
// @Summary 获取评论回复列表
// @Description 获取指定评论的回复列表，支持分页
// @Tags 评论
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageData}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /comments/{id}/replies [get]
func (c *CommentController) GetCommentReplies(ctx *gin.Context) {
	// 获取评论ID
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的评论ID")
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

	// 获取当前用户ID（可选）
	var userID *uint
	if uid, exists := ctx.Get("user_id"); exists {
		temp := uid.(uint)
		userID = &temp
	}

	// 获取评论回复列表
	replies, total, err := c.commentService.GetCommentReplies(uint(commentID), page, size, userID)
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, "获取回复列表失败")
		return
	}

	response.SuccessPage(ctx, replies, total, page, size)
}

// GetUserComments 获取用户评论列表
// @Summary 获取用户评论列表
// @Description 获取当前用户的评论列表，支持分页
// @Tags 评论
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageData}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /comments/mine [get]
// @Security BearerAuth
func (c *CommentController) GetUserComments(ctx *gin.Context) {
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

	// 获取用户评论列表
	comments, total, err := c.commentService.GetUserComments(userID.(uint), page, size)
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, "获取评论列表失败")
		return
	}

	response.SuccessPage(ctx, comments, total, page, size)
}

// ApproveComment 审核通过评论（管理员功能）
// @Summary 审核通过评论
// @Description 管理员审核通过指定评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /comments/{id}/approve [post]
// @Security BearerAuth
func (c *CommentController) ApproveComment(ctx *gin.Context) {
	// 获取评论ID
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的评论ID")
		return
	}

	// 审核通过评论
	err = c.commentService.ApproveComment(uint(commentID))
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	response.SuccessWithMessage(ctx, nil, "评论审核通过")
}

// RejectComment 拒绝评论（管理员功能）
// @Summary 拒绝评论
// @Description 管理员拒绝指定评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /comments/{id}/reject [post]
// @Security BearerAuth
func (c *CommentController) RejectComment(ctx *gin.Context) {
	// 获取评论ID
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的评论ID")
		return
	}

	// 拒绝评论
	err = c.commentService.RejectComment(uint(commentID))
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.ErrorWithMessage(ctx, appErr.Code, appErr.Message)
		} else {
			response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		}
		return
	}

	response.SuccessWithMessage(ctx, nil, "评论已拒绝")
}