package middleware

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// CORS 中间件，处理跨域请求
func CORS() gin.HandlerFunc {
	_ = godotenv.Load() // 加载 .env 文件，忽略错误

	origins := os.Getenv("CORS_ALLOW_ORIGINS")
	var allowOrigins []string
	if origins == "" || origins == "*" {
		allowOrigins = []string{"*"}
	} else {
		allowOrigins = strings.Split(origins, ",")
	}

	config := cors.DefaultConfig()
	if len(allowOrigins) == 1 && allowOrigins[0] == "*" {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = allowOrigins
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept"}

	return cors.New(config)
}
