package middlewares

import (
	"blog-api/pkg/response"
	"github.com/gin-gonic/gin"
)

// AdminMiddleware 管理员权限中间件
// 验证用户是否具有管理员权限
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先确保用户已登录（需要在AdminMiddleware之前使用AuthMiddleware）
		_, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithMessage(c, response.UNAUTHORIZED, "用户未登录")
			c.Abort()
			return
		}

		// 获取用户角色
		userRole, roleExists := c.Get("user_role")
		if !roleExists {
			response.ErrorWithMessage(c, response.FORBIDDEN, "无法获取用户角色信息")
			c.Abort()
			return
		}

		// 检查是否为管理员
		if userRole.(string) != "admin" {
			response.ErrorWithMessage(c, response.FORBIDDEN, "需要管理员权限")
			c.Abort()
			return
		}

		// 通过权限验证，继续处理请求
		c.Next()
	}
}

// OptionalAdminMiddleware 可选的管理员权限中间件
// 如果用户已登录且是管理员，则设置相关标识，但不会阻止非管理员用户访问
func OptionalAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查用户是否已登录
		_, exists := c.Get("user_id")
		if !exists {
			c.Set("is_admin", false)
			c.Next()
			return
		}

		// 获取用户角色
		userRole, roleExists := c.Get("user_role")
		if !roleExists {
			c.Set("is_admin", false)
			c.Next()
			return
		}

		// 设置管理员标识
		isAdmin := userRole.(string) == "admin"
		c.Set("is_admin", isAdmin)
		
		c.Next()
	}
}