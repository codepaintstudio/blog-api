package middlewares

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders 安全头中间件
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// X-Content-Type-Options: 防止MIME类型混淆攻击
		c.Header("X-Content-Type-Options", "nosniff")
		
		// X-Frame-Options: 防止点击劫持攻击
		c.Header("X-Frame-Options", "DENY")
		
		// X-XSS-Protection: 启用浏览器XSS防护
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Referrer-Policy: 控制Referrer信息泄露
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Content-Security-Policy: 内容安全策略
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self'; " +
			"connect-src 'self'; " +
			"media-src 'self'; " +
			"object-src 'none'; " +
			"child-src 'none'; " +
			"frame-ancestors 'none'; " +
			"form-action 'self'; " +
			"base-uri 'self'"
		c.Header("Content-Security-Policy", csp)
		
		// Strict-Transport-Security: 强制HTTPS (仅在HTTPS环境下设置)
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}
		
		// Permissions-Policy: 控制浏览器功能权限
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		c.Next()
	}
}