package controllers

import (
	"strconv"

	"blog-api/internal/middlewares"
	"blog-api/internal/services"
	"blog-api/pkg/logger"
	"blog-api/pkg/response"
	
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ArticleController 文章控制器
type ArticleController struct {
	articleService services.ArticleService
	validator      *validator.Validate
}

// NewArticleController 创建新的文章控制器
func NewArticleController(articleService services.ArticleService) *ArticleController {
	return &ArticleController{
		articleService: articleService,
		validator:      validator.New(),
	}
}

// Create 创建文章
// @Summary 创建文章
// @Description 创建新文章
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body services.CreateArticleRequest true "文章信息"
// @Success 200 {object} response.Response{data=services.ArticleResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /articles [post]
func (ac *ArticleController) Create(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := middlewares.GetCurrentUserID(c)
	if !exists {
		response.Error(c, response.UNAUTHORIZED)
		return
	}
	
	var req services.CreateArticleRequest
	
	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("创建文章参数绑定失败",
			logger.Err("error", err),
			logger.Int("user_id", int(userID)),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "参数格式错误")
		return
	}
	
	// 验证参数
	if err := ac.validator.Struct(&req); err != nil {
		logger.Warn("创建文章参数验证失败",
			logger.Err("error", err),
			logger.Int("user_id", int(userID)),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "参数验证失败")
		return
	}
	
	// 调用服务层创建文章
	article, err := ac.articleService.Create(userID, &req)
	if err != nil {
		logger.Error("创建文章失败",
			logger.Err("error", err),
			logger.Int("user_id", int(userID)),
		)
		
		if err.Error() == "分类不存在" {
			response.ErrorWithMessage(c, response.NOT_FOUND, err.Error())
		} else {
			response.ErrorWithMessage(c, response.ERROR, "创建文章失败")
		}
		return
	}
	
	response.SuccessWithMessage(c, article, "文章创建成功")
}

// GetByID 获取文章详情
// @Summary 获取文章详情
// @Description 根据ID获取文章详情
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response{data=services.ArticleDetailResponse}
// @Failure 404 {object} response.Response
// @Router /articles/{id} [get]
func (ac *ArticleController) GetByID(c *gin.Context) {
	// 获取文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "无效的文章ID")
		return
	}
	
	// 获取当前用户ID（可选）
	userID, _ := middlewares.GetCurrentUserID(c)
	var currentUserID *uint
	if userID != 0 {
		currentUserID = &userID
	}
	
	// 调用服务层获取文章
	article, err := ac.articleService.GetByID(uint(id), currentUserID)
	if err != nil {
		logger.Warn("获取文章失败",
			logger.Err("error", err),
			logger.String("article_id", idStr),
		)
		
		if err.Error() == "文章不存在" {
			response.Error(c, response.NOT_FOUND)
		} else {
			response.ErrorWithMessage(c, response.ERROR, "获取文章失败")
		}
		return
	}
	
	response.Success(c, article)
}

// Update 更新文章
// @Summary 更新文章
// @Description 更新文章信息
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Param body body services.UpdateArticleRequest true "更新信息"
// @Success 200 {object} response.Response{data=services.ArticleResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /articles/{id} [put]
func (ac *ArticleController) Update(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := middlewares.GetCurrentUserID(c)
	if !exists {
		response.Error(c, response.UNAUTHORIZED)
		return
	}
	
	// 获取文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "无效的文章ID")
		return
	}
	
	var req services.UpdateArticleRequest
	
	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("更新文章参数绑定失败",
			logger.Err("error", err),
			logger.Int("user_id", int(userID)),
			logger.String("article_id", idStr),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "参数格式错误")
		return
	}
	
	// 验证参数
	if err := ac.validator.Struct(&req); err != nil {
		logger.Warn("更新文章参数验证失败",
			logger.Err("error", err),
			logger.Int("user_id", int(userID)),
			logger.String("article_id", idStr),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "参数验证失败")
		return
	}
	
	// 调用服务层更新文章
	article, err := ac.articleService.Update(uint(id), userID, &req)
	if err != nil {
		logger.Error("更新文章失败",
			logger.Err("error", err),
			logger.Int("user_id", int(userID)),
			logger.String("article_id", idStr),
		)
		
		switch err.Error() {
		case "文章不存在":
			response.Error(c, response.NOT_FOUND)
		case "无权限操作":
			response.Error(c, response.FORBIDDEN)
		case "分类不存在":
			response.ErrorWithMessage(c, response.NOT_FOUND, err.Error())
		default:
			response.ErrorWithMessage(c, response.ERROR, "更新文章失败")
		}
		return
	}
	
	response.SuccessWithMessage(c, article, "文章更新成功")
}

// Delete 删除文章
// @Summary 删除文章
// @Description 删除文章
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /articles/{id} [delete]
func (ac *ArticleController) Delete(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := middlewares.GetCurrentUserID(c)
	if !exists {
		response.Error(c, response.UNAUTHORIZED)
		return
	}
	
	// 获取文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "无效的文章ID")
		return
	}
	
	// 调用服务层删除文章
	if err := ac.articleService.Delete(uint(id), userID); err != nil {
		logger.Error("删除文章失败",
			logger.Err("error", err),
			logger.Int("user_id", int(userID)),
			logger.String("article_id", idStr),
		)
		
		switch err.Error() {
		case "文章不存在":
			response.Error(c, response.NOT_FOUND)
		case "无权限操作":
			response.Error(c, response.FORBIDDEN)
		default:
			response.ErrorWithMessage(c, response.ERROR, "删除文章失败")
		}
		return
	}
	
	response.SuccessWithMessage(c, nil, "文章删除成功")
}

// List 获取文章列表
// @Summary 获取文章列表
// @Description 获取文章列表（支持分页和筛选）
// @Tags 文章
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param category_id query int false "分类ID"
// @Param sort_by query string false "排序字段" Enums(created_at, updated_at, view_count, like_count)
// @Param sort_order query string false "排序方向" Enums(asc, desc)
// @Success 200 {object} response.Response{data=services.ListArticleResponse}
// @Router /articles [get]
func (ac *ArticleController) List(c *gin.Context) {
	var req services.ListArticleRequest
	
	// 绑定查询参数
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warn("获取文章列表参数绑定失败",
			logger.Err("error", err),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "参数格式错误")
		return
	}
	
	// 调用服务层获取文章列表
	articles, err := ac.articleService.List(&req)
	if err != nil {
		logger.Error("获取文章列表失败",
			logger.Err("error", err),
		)
		response.ErrorWithMessage(c, response.ERROR, "获取文章列表失败")
		return
	}
	
	response.Success(c, articles)
}

// Search 搜索文章
// @Summary 搜索文章
// @Description 根据关键词搜索文章
// @Tags 文章
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=services.ListArticleResponse}
// @Router /articles/search [get]
func (ac *ArticleController) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "搜索关键词不能为空")
		return
	}
	
	// 获取分页参数
	page := 1
	pageSize := 10
	
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}
	
	// 调用服务层搜索文章
	articles, err := ac.articleService.Search(keyword, page, pageSize)
	if err != nil {
		logger.Error("搜索文章失败",
			logger.Err("error", err),
			logger.String("keyword", keyword),
		)
		response.ErrorWithMessage(c, response.ERROR, "搜索失败")
		return
	}
	
	response.Success(c, articles)
}