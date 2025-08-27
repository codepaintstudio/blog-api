package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`           // 业务状态码
	Message string      `json:"message"`        // 响应消息
	Data    interface{} `json:"data,omitempty"` // 响应数据
}

// PageData 分页数据结构
type PageData struct {
	List     interface{} `json:"list"`      // 数据列表
	Total    int64       `json:"total"`     // 总数
	Page     int         `json:"page"`      // 当前页
	PageSize int         `json:"page_size"` // 每页大小
	Pages    int         `json:"pages"`     // 总页数
}

// 业务状态码常量
const (
	SUCCESS = 200
	ERROR   = 500

	// 认证相关
	UNAUTHORIZED = 401
	FORBIDDEN    = 403

	// 参数相关
	INVALID_PARAMS = 400

	// 数据相关
	NOT_FOUND      = 404
	ALREADY_EXISTS = 409
)

// 响应消息映射
var codeMessages = map[int]string{
	SUCCESS:        "操作成功",
	ERROR:          "操作失败",
	UNAUTHORIZED:   "未授权",
	FORBIDDEN:      "禁止访问",
	INVALID_PARAMS: "参数错误",
	NOT_FOUND:      "数据不存在",
	ALREADY_EXISTS: "数据已存在",
}

// getCodeMessage 获取状态码对应的消息
func getCodeMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    SUCCESS,
		Message: getCodeMessage(SUCCESS),
		Data:    data,
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    SUCCESS,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: getCodeMessage(code),
	})
}

// ErrorWithMessage 带自定义消息的错误响应
func ErrorWithMessage(c *gin.Context, code int, message string) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}

	pageData := PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}

	Success(c, pageData)
}

// getHTTPStatus 根据业务状态码获取HTTP状态码
func getHTTPStatus(code int) int {
	switch code {
	case SUCCESS:
		return http.StatusOK
	case UNAUTHORIZED:
		return http.StatusUnauthorized
	case FORBIDDEN:
		return http.StatusForbidden
	case INVALID_PARAMS:
		return http.StatusBadRequest
	case NOT_FOUND:
		return http.StatusNotFound
	case ALREADY_EXISTS:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}