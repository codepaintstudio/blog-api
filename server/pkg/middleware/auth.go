package middleware

import (
	"server/pkg/response"
	"server/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// validateToken 验证token并返回claims
func validateToken(c *gin.Context) (*utils.JwtClaims, bool) {
	// 获取 Authorization Header
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		c.JSON(401, response.Error(response.StatusUnauthorized, "Authorization header is required"))
		c.Abort()
		return nil, false
	}

	// 解析 Bearer Token
	parts := strings.Split(authorization, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(401, response.Error(response.StatusUnauthorized, "Invalid Authorization header format"))
		c.Abort()
		return nil, false
	}

	// 解析并验证Token
	claims, err := utils.ParseToken(parts[1])
	if err != nil {
		c.JSON(401, response.Error(response.StatusUnauthorized, "Invalid token"))
		c.Abort()
		return nil, false
	}

	return claims, true
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := validateToken(c)
		if !ok {
			return
		}

		// 将用户ID存入上下文
		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
