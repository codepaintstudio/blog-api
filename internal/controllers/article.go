package controllers

import (
	"server/internal/models"
	"server/internal/services"
	"server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleController struct {
	articleService *services.ArticleService
}

func NewArticleController(db *gorm.DB) *ArticleController {
	return &ArticleController{
		articleService: services.NewArticleService(db),
	}
}

// Create 创建帖子
func (c *ArticleController) Create(ctx *gin.Context) {
	var request models.CreateArticleRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, response.Error(response.StatusBadRequest, "Invalid request format"))
		return
	}

	userId := ctx.GetInt("user_id")
	article, err := c.articleService.Create(userId, &request)
	if err != nil {
		ctx.JSON(500, response.Error(response.StatusInternalError, err.Error()))
		return
	}

	ctx.JSON(200, response.SuccessWithMessage("Create article successfully", article))
}

// Update 更新帖子
func (c *ArticleController) Update(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, response.Error(response.StatusBadRequest, "Invalid article ID"))
		return
	}

	var request models.UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, response.Error(response.StatusBadRequest, "Invalid request format"))
		return
	}

	userId := ctx.GetInt("user_id")
	article, err := c.articleService.Update(userId, articleId, &request)
	if err != nil {
		if err.Error() == "unauthorized to update this article" {
			ctx.JSON(403, response.Error(response.StatusForbidden, err.Error()))
			return
		}
		ctx.JSON(500, response.Error(response.StatusInternalError, err.Error()))
		return
	}

	ctx.JSON(200, response.SuccessWithMessage("Update article successfully", article))
}

// Delete 删除帖子
func (c *ArticleController) Delete(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, response.Error(response.StatusBadRequest, "Invalid article ID"))
		return
	}

	userId := ctx.GetInt("user_id")
	err = c.articleService.Delete(userId, articleId)
	if err != nil {
		if err.Error() == "unauthorized to delete this article" {
			ctx.JSON(403, response.Error(response.StatusForbidden, err.Error()))
			return
		}
		ctx.JSON(500, response.Error(response.StatusInternalError, err.Error()))
		return
	}

	ctx.JSON(200, response.SuccessWithMessage("Delete article successfully", nil))
}

// GetById 获取帖子详情
func (c *ArticleController) GetById(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, response.Error(response.StatusBadRequest, "Invalid article ID"))
		return
	}

	article, err := c.articleService.GetById(articleId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(404, response.Error(response.StatusNotFound, "Article not found"))
			return
		}
		ctx.JSON(500, response.Error(response.StatusInternalError, err.Error()))
		return
	}

	ctx.JSON(200, response.SuccessWithMessage("Get article successfully", article))
}

// List 获取帖子列表
func (c *ArticleController) List(ctx *gin.Context) {
	var request models.ArticleListRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(400, response.Error(response.StatusBadRequest, "Invalid request format"))
		return
	}

	// 设置默认值
	if request.Page <= 0 {
		request.Page = 1
	}
	if request.Size <= 0 {
		request.Size = 10
	}

	data, err := c.articleService.List(&request)
	if err != nil {
		ctx.JSON(500, response.Error(response.StatusInternalError, err.Error()))
		return
	}

	ctx.JSON(200, response.SuccessWithMessage("Get articles successfully", data))
}

// GetStats 获取文章统计信息
func (c *ArticleController) GetStats(ctx *gin.Context) {
	stats, err := c.articleService.GetStats()
	if err != nil {
		ctx.JSON(500, response.Error(response.StatusInternalError, err.Error()))
		return
	}

	ctx.JSON(200, response.SuccessWithMessage("Get stats successfully", stats))
}
