package controllers

import (
	"blog-api/internal/services"
	"blog-api/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FileController 文件控制器
type FileController struct {
	fileService services.FileService
}

// NewFileController 创建文件控制器实例
func NewFileController(fileService services.FileService) *FileController {
	return &FileController{
		fileService: fileService,
	}
}

// UploadRequest 文件上传请求
type UploadRequest struct {
	File string `form:"file" binding:"required"`
}

// FileResponse 文件响应
type FileResponse struct {
	ID           uint   `json:"id"`
	Filename     string `json:"filename"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mime_type"`
	URL          string `json:"url"`
	Hash         string `json:"hash,omitempty"`      // 可选返回哈希值
	RefCount     int    `json:"ref_count,omitempty"` // 可选返回引用计数
	CreatedAt    string `json:"created_at"`
}

// Upload 上传文件
// @Summary 上传文件
// @Description 上传图片文件，支持去重机制
// @Tags 文件
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "上传的文件"
// @Success 200 {object} response.Response{data=FileResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /files/upload [post]
func (c *FileController) Upload(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, response.UNAUTHORIZED)
		return
	}

	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "请选择要上传的文件")
		return
	}

	// 调用服务层上传文件
	fileModel, err := c.fileService.Upload(userID.(uint), file)
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	// 返回文件信息
	fileResp := &FileResponse{
		ID:           fileModel.ID,
		Filename:     fileModel.Filename,
		OriginalName: fileModel.OriginalName,
		Size:         fileModel.Size,
		MimeType:     fileModel.MimeType,
		URL:          fileModel.URL,
		CreatedAt:    fileModel.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	response.Success(ctx, fileResp)
}

// GetFile 获取文件信息
// @Summary 获取文件信息
// @Description 根据文件ID获取文件详细信息
// @Tags 文件
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} response.Response{data=FileResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /files/{id} [get]
func (c *FileController) GetFile(ctx *gin.Context) {
	// 获取文件ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的文件ID")
		return
	}

	// 获取文件信息
	file, err := c.fileService.GetByID(uint(id))
	if err != nil {
		response.ErrorWithMessage(ctx, response.NOT_FOUND, err.Error())
		return
	}

	// 返回文件信息
	fileResp := &FileResponse{
		ID:           file.ID,
		Filename:     file.Filename,
		OriginalName: file.OriginalName,
		Size:         file.Size,
		MimeType:     file.MimeType,
		URL:          file.URL,
		Hash:         file.Hash,
		RefCount:     file.RefCount,
		CreatedAt:    file.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	response.Success(ctx, fileResp)
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Description 删除文件（减少引用计数）
// @Tags 文件
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /files/{id} [delete]
func (c *FileController) DeleteFile(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, response.UNAUTHORIZED)
		return
	}

	// 获取文件ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.BAD_REQUEST, "无效的文件ID")
		return
	}

	// 删除文件
	if err := c.fileService.Delete(uint(id), userID.(uint)); err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, nil, "文件删除成功")
}

// GetUserFiles 获取用户文件列表
// @Summary 获取用户文件列表
// @Description 分页获取当前用户的文件列表
// @Tags 文件
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{items=[]FileResponse}}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /files [get]
func (c *FileController) GetUserFiles(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, response.UNAUTHORIZED)
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))

	// 获取用户文件列表
	files, total, err := c.fileService.GetUserFiles(userID.(uint), page, size)
	if err != nil {
		response.ErrorWithMessage(ctx, response.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	// 转换为响应格式
	var fileResponses []FileResponse
	for _, file := range files {
		fileResponses = append(fileResponses, FileResponse{
			ID:           file.ID,
			Filename:     file.Filename,
			OriginalName: file.OriginalName,
			Size:         file.Size,
			MimeType:     file.MimeType,
			URL:          file.URL,
			CreatedAt:    file.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// 返回分页数据
	pageData := response.PageData{
		List:     fileResponses,
		Total:    total,
		Page:     page,
		PageSize: size,
	}

	response.Success(ctx, pageData)
}

// ServeFile 提供静态文件服务
// @Summary 访问静态文件
// @Description 通过文件名访问上传的静态文件
// @Tags 文件
// @Produce octet-stream
// @Param filename path string true "文件名"
// @Success 200 {file} file
// @Failure 404 {object} response.Response
// @Router /uploads/{filename} [get]
func (c *FileController) ServeFile(ctx *gin.Context) {
	filename := ctx.Param("filename")

	// 这个方法将在路由中直接使用gin.Static处理
	// 这里只是用于Swagger文档生成
	ctx.File("./uploads/" + filename)
}
