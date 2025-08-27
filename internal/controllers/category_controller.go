package controllers

import (
	"strconv"

	"blog-api/internal/services"
	"blog-api/pkg/errors"
	"blog-api/pkg/logger"
	"blog-api/pkg/response"

	"github.com/gin-gonic/gin"
)

// CategoryController 分类控制器
type CategoryController struct {
	categoryService services.CategoryService
}

// NewCategoryController 创建新的分类控制器
func NewCategoryController(categoryService services.CategoryService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

// Create 创建分类
// @Summary 创建分类
// @Description 创建新的文章分类
// @Tags 分类
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category body services.CreateCategoryRequest true "分类信息"
// @Success 201 {object} response.Response{data=services.CategoryResponse} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /categories [post]
func (cc *CategoryController) Create(c *gin.Context) {
	var req services.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("分类创建参数绑定失败",
			logger.Err("error", err),
			logger.String("ip", c.ClientIP()),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, errors.ErrInvalidParams.Error())
		return
	}

	category, err := cc.categoryService.Create(&req)
	if err != nil {
		logger.Error("分类创建失败",
			logger.Err("error", err),
			logger.String("name", req.Name),
		)

		switch err.Error() {
		case "分类名称已存在":
			response.ErrorWithMessage(c, response.INVALID_PARAMS, err.Error())
		default:
			response.ErrorWithMessage(c, response.ERROR, "创建分类失败")
		}
		return
	}

	// 转换为响应格式
	categoryResp := &services.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Color:       category.Color,
		Icon:        category.Icon,
		Sort:        category.Sort,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	response.SuccessWithMessage(c, categoryResp, "分类创建成功")
}

// GetByID 根据ID获取分类详情
// @Summary 获取分类详情
// @Description 根据分类ID获取分类详细信息
// @Tags 分类
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response{data=services.CategoryResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "分类不存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /categories/{id} [get]
func (cc *CategoryController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "无效的分类ID")
		return
	}

	category, err := cc.categoryService.GetByID(uint(id))
	if err != nil {
		logger.Error("获取分类详情失败",
			logger.Err("error", err),
			logger.Int("category_id", int(id)),
		)

		switch err.Error() {
		case "分类不存在":
			response.ErrorWithMessage(c, response.NOT_FOUND, err.Error())
		default:
			response.ErrorWithMessage(c, response.ERROR, "获取分类详情失败")
		}
		return
	}

	// 转换为响应格式
	categoryResp := &services.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Color:       category.Color,
		Icon:        category.Icon,
		Sort:        category.Sort,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	response.SuccessWithMessage(c, categoryResp, "获取分类详情成功")
}

// Update 更新分类
// @Summary 更新分类
// @Description 更新分类信息
// @Tags 分类
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Param category body services.UpdateCategoryRequest true "分类信息"
// @Success 200 {object} response.Response{data=services.CategoryResponse} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "分类不存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /categories/{id} [put]
func (cc *CategoryController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "无效的分类ID")
		return
	}

	var req services.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("分类更新参数绑定失败",
			logger.Err("error", err),
			logger.String("ip", c.ClientIP()),
		)
		response.ErrorWithMessage(c, response.INVALID_PARAMS, errors.ErrInvalidParams.Error())
		return
	}

	category, err := cc.categoryService.Update(uint(id), &req)
	if err != nil {
		logger.Error("分类更新失败",
			logger.Err("error", err),
			logger.Int("category_id", int(id)),
		)

		switch err.Error() {
		case "分类不存在":
			response.ErrorWithMessage(c, response.NOT_FOUND, err.Error())
		case "分类名称已存在":
			response.ErrorWithMessage(c, response.INVALID_PARAMS, err.Error())
		default:
			response.ErrorWithMessage(c, response.ERROR, "更新分类失败")
		}
		return
	}

	// 转换为响应格式
	categoryResp := &services.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Color:       category.Color,
		Icon:        category.Icon,
		Sort:        category.Sort,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	response.SuccessWithMessage(c, categoryResp, "分类更新成功")
}

// Delete 删除分类
// @Summary 删除分类
// @Description 删除指定的分类
// @Tags 分类
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "分类不存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /categories/{id} [delete]
func (cc *CategoryController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(c, response.INVALID_PARAMS, "无效的分类ID")
		return
	}

	err = cc.categoryService.Delete(uint(id))
	if err != nil {
		logger.Error("分类删除失败",
			logger.Err("error", err),
			logger.Int("category_id", int(id)),
		)

		switch err.Error() {
		case "分类不存在":
			response.ErrorWithMessage(c, response.NOT_FOUND, err.Error())
		case "该分类下还有文章，无法删除":
			response.ErrorWithMessage(c, response.INVALID_PARAMS, err.Error())
		default:
			response.ErrorWithMessage(c, response.ERROR, "删除分类失败")
		}
		return
	}

	response.SuccessWithMessage(c, nil, "分类删除成功")
}

// List 获取所有分类列表
// @Summary 获取分类列表
// @Description 获取所有分类列表，包含文章数量统计
// @Tags 分类
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]services.CategoryListResponse} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /categories [get]
func (cc *CategoryController) List(c *gin.Context) {
	categories, err := cc.categoryService.List()
	if err != nil {
		logger.Error("获取分类列表失败", logger.Err("error", err))
		response.ErrorWithMessage(c, response.ERROR, "获取分类列表失败")
		return
	}

	response.SuccessWithMessage(c, categories, "获取分类列表成功")
}

// ListActive 获取活跃分类列表
// @Summary 获取活跃分类列表
// @Description 获取所有活跃状态的分类列表
// @Tags 分类
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]services.CategoryListResponse} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /categories/active [get]
func (cc *CategoryController) ListActive(c *gin.Context) {
	categories, err := cc.categoryService.ListActive()
	if err != nil {
		logger.Error("获取活跃分类列表失败", logger.Err("error", err))
		response.ErrorWithMessage(c, response.ERROR, "获取活跃分类列表失败")
		return
	}

	response.SuccessWithMessage(c, categories, "获取活跃分类列表成功")
}
