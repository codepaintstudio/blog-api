package middlewares

import (
	"bytes"
	"io"
	"time"

	"blog-api/pkg/logger"
	
	"github.com/gin-gonic/gin"
)

// responseWriter 自定义响应写入器，用于捕获响应内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggerMiddleware 自定义日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		
		// 读取请求体（如果需要）
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
		
		// 创建自定义响应写入器
		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:          bytes.NewBufferString(""),
		}
		c.Writer = blw
		
		// 处理请求
		c.Next()
		
		// 结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		
		// 获取请求信息
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()
		
		// 记录日志
		logFields := []logger.Field{
			logger.String("method", method),
			logger.String("path", path),
			logger.Int("status", statusCode),
			logger.String("ip", clientIP),
			logger.String("user_agent", userAgent),
			logger.String("latency", latency.String()),
		}
		
		// 根据状态码选择日志级别
		logMessage := "HTTP Request"
		if statusCode >= 500 {
			logger.Error(logMessage, logFields...)
		} else if statusCode >= 400 {
			logger.Warn(logMessage, logFields...)
		} else {
			logger.Info(logMessage, logFields...)
		}
	}
}