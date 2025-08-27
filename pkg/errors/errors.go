package errors

import (
	"fmt"
)

// AppError 应用错误类型
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Code: %d, Message: %s, Error: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// NewAppError 创建新的应用错误
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 预定义错误
var (
	// 通用错误
	ErrInternalServer = &AppError{Code: 500, Message: "服务器内部错误"}
	ErrInvalidParams  = &AppError{Code: 400, Message: "参数错误"}
	ErrNotFound       = &AppError{Code: 404, Message: "资源不存在"}
	ErrConflict       = &AppError{Code: 409, Message: "资源冲突"}
	
	// 认证授权错误
	ErrUnauthorized     = &AppError{Code: 401, Message: "未授权"}
	ErrForbidden        = &AppError{Code: 403, Message: "禁止访问"}
	ErrInvalidToken     = &AppError{Code: 401, Message: "无效的令牌"}
	ErrTokenExpired     = &AppError{Code: 401, Message: "令牌已过期"}
	ErrInvalidPassword  = &AppError{Code: 401, Message: "密码错误"}
	
	// 用户相关错误
	ErrUserNotFound     = &AppError{Code: 404, Message: "用户不存在"}
	ErrUserExists       = &AppError{Code: 409, Message: "用户已存在"}
	ErrEmailExists      = &AppError{Code: 409, Message: "邮箱已被注册"}
	ErrUsernameExists   = &AppError{Code: 409, Message: "用户名已被注册"}
	ErrUserInactive     = &AppError{Code: 403, Message: "用户账号已被禁用"}
	
	// 文章相关错误
	ErrArticleNotFound  = &AppError{Code: 404, Message: "文章不存在"}
	ErrArticleExists    = &AppError{Code: 409, Message: "文章已存在"}
	ErrNoPermission     = &AppError{Code: 403, Message: "无操作权限"}
	
	// 分类相关错误
	ErrCategoryNotFound = &AppError{Code: 404, Message: "分类不存在"}
	ErrCategoryExists   = &AppError{Code: 409, Message: "分类已存在"}
	
	// 评论相关错误
	ErrCommentNotFound  = &AppError{Code: 404, Message: "评论不存在"}
	
	// 文件相关错误
	ErrFileNotFound     = &AppError{Code: 404, Message: "文件不存在"}
	ErrFileTooBig       = &AppError{Code: 413, Message: "文件过大"}
	ErrFileTypeNotAllowed = &AppError{Code: 415, Message: "不支持的文件类型"}
	
	// 限流错误
	ErrTooManyRequests  = &AppError{Code: 429, Message: "请求过于频繁"}
	
	// 数据库相关错误
	ErrDatabase         = &AppError{Code: 500, Message: "数据库错误"}
	ErrRedis            = &AppError{Code: 500, Message: "缓存服务错误"}
)

// IsAppError 检查是否为应用错误
func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}

// WrapError 包装错误
func WrapError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}