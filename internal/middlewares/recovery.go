package middlewares

import (
	"fmt"
	"net"
	"os"
	"runtime/debug"
	"strings"

	"blog-api/pkg/logger"
	"blog-api/pkg/response"
	
	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware 自定义错误恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic错误和堆栈信息
				stack := debug.Stack()
				
				logger.Error("Panic recovered",
					logger.String("error", fmt.Sprintf("%v", err)),
					logger.String("stack", string(stack)),
					logger.String("path", c.Request.URL.Path),
					logger.String("method", c.Request.Method),
					logger.String("ip", c.ClientIP()),
				)
				
				// 检查连接是否已断开
				if checkBrokenPipe(err) {
					c.Error(fmt.Errorf("%s", err))
					c.Abort()
					return
				}
				
				// 返回500错误响应
				response.ErrorWithMessage(c, response.ERROR, "服务器内部错误")
				c.Abort()
			}
		}()
		
		c.Next()
	}
}

// checkBrokenPipe 检查是否为连接断开错误
func checkBrokenPipe(err interface{}) bool {
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
				strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
				return true
			}
		}
	}
	return false
}