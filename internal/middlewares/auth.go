package middlewares

import (
	"strings"

	"blog-api/internal/utils"
	"blog-api/pkg/logger"
	"blog-api/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	// AuthHeader 认证头名称
	AuthHeader = "Authorization"
	// BearerPrefix Bearer前缀
	BearerPrefix = "Bearer "
	// UserIDKey 用户ID在上下文中的键名
	UserIDKey = "user_id"
	// UsernameKey 用户名在上下文中的键名
	UsernameKey = "username"
	// UserRoleKey 用户角色在上下文中的键名
	UserRoleKey = "user_role"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthHeader)
		if authHeader == "" {
			logger.Warn("缺少Authorization头",
				logger.String("path", c.Request.URL.Path),
				logger.String("ip", c.ClientIP()),
			)
			response.Error(c, response.UNAUTHORIZED)
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			logger.Warn("无效的Authorization头格式",
				logger.String("header", authHeader),
				logger.String("path", c.Request.URL.Path),
				logger.String("ip", c.ClientIP()),
			)
			response.ErrorWithMessage(c, response.UNAUTHORIZED, "无效的token格式")
			c.Abort()
			return
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			response.ErrorWithMessage(c, response.UNAUTHORIZED, "token不能为空")
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			var message string
			switch err {
			case utils.ErrTokenExpired:
				message = "token已过期"
			case utils.ErrTokenNotValidYet:
				message = "token尚未生效"
			case utils.ErrTokenMalformed:
				message = "token格式错误"
			default:
				message = "无效的token"
			}

			logger.Warn("token验证失败",
				logger.String("error", err.Error()),
				logger.String("path", c.Request.URL.Path),
				logger.String("ip", c.ClientIP()),
			)

			response.ErrorWithMessage(c, response.UNAUTHORIZED, message)
			c.Abort()
			return
		}

		// 验证token类型
		if claims.Type != "access" {
			response.ErrorWithMessage(c, response.UNAUTHORIZED, "无效的token类型")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set(UserIDKey, claims.UserID)
		c.Set(UsernameKey, claims.Username)
		c.Set(UserRoleKey, claims.Role)

		logger.Debug("用户认证成功",
			logger.Int("user_id", int(claims.UserID)),
			logger.String("username", claims.Username),
			logger.String("role", claims.Role),
		)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（用户可以选择是否登录）
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthHeader)
		if authHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			c.Next()
			return
		}

		// 尝试解析token，失败则继续，不阻止请求
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		if claims.Type == "access" {
			c.Set(UserIDKey, claims.UserID)
			c.Set(UsernameKey, claims.Username)
			c.Set(UserRoleKey, claims.Role)
		}

		c.Next()
	}
}

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthHeader)
		if authHeader == "" {
			logger.Warn("缺少Authorization头",
				logger.String("path", c.Request.URL.Path),
				logger.String("ip", c.ClientIP()),
			)
			response.Error(c, response.UNAUTHORIZED)
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			logger.Warn("无效的Authorization头格式",
				logger.String("header", authHeader),
				logger.String("path", c.Request.URL.Path),
				logger.String("ip", c.ClientIP()),
			)
			response.ErrorWithMessage(c, response.UNAUTHORIZED, "无效的token格式")
			c.Abort()
			return
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			response.ErrorWithMessage(c, response.UNAUTHORIZED, "token不能为空")
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			var message string
			switch err {
			case utils.ErrTokenExpired:
				message = "token已过期"
			case utils.ErrTokenNotValidYet:
				message = "token尚未生效"
			case utils.ErrTokenMalformed:
				message = "token格式错误"
			default:
				message = "无效的token"
			}

			logger.Warn("token验证失败",
				logger.String("error", err.Error()),
				logger.String("path", c.Request.URL.Path),
				logger.String("ip", c.ClientIP()),
			)

			response.ErrorWithMessage(c, response.UNAUTHORIZED, message)
			c.Abort()
			return
		}

		// 验证token类型
		if claims.Type != "access" {
			response.ErrorWithMessage(c, response.UNAUTHORIZED, "无效的token类型")
			c.Abort()
			return
		}

		// 检查用户角色（在设置上下文之前检查）
		if claims.Role != "admin" {
			logger.Warn("非管理员尝试访问管理员接口",
				logger.String("role", claims.Role),
				logger.String("path", c.Request.URL.Path),
				logger.String("ip", c.ClientIP()),
			)
			response.Error(c, response.FORBIDDEN)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set(UserIDKey, claims.UserID)
		c.Set(UsernameKey, claims.Username)
		c.Set(UserRoleKey, claims.Role)

		logger.Debug("管理员认证成功",
			logger.Int("user_id", int(claims.UserID)),
			logger.String("username", claims.Username),
			logger.String("role", claims.Role),
		)

		c.Next()
	}
}

// GetCurrentUserID 从上下文中获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	return id, ok
}

// GetCurrentUsername 从上下文中获取当前用户名
func GetCurrentUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get(UsernameKey)
	if !exists {
		return "", false
	}

	name, ok := username.(string)
	return name, ok
}

// GetCurrentUserRole 从上下文中获取当前用户角色
func GetCurrentUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(UserRoleKey)
	if !exists {
		return "", false
	}

	r, ok := role.(string)
	return r, ok
}
