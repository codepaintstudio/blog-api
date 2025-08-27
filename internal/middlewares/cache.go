package middlewares

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"blog-api/pkg/cache"
	"blog-api/pkg/logger"
	"blog-api/pkg/response"

	"github.com/gin-gonic/gin"
)

// CacheMiddleware 缓存中间件配置
type CacheMiddleware struct {
	strategy *cache.CacheStrategy
	config   *CacheConfig
}

// CacheConfig 缓存配置
type CacheConfig struct {
	// 需要缓存的路径配置
	CachePaths map[string]CachePathConfig
	// 默认缓存时间
	DefaultTTL time.Duration
	// 是否启用缓存
	Enabled bool
}

// CachePathConfig 路径缓存配置
type CachePathConfig struct {
	TTL     time.Duration // 缓存时间
	Methods []string      // 缓存的HTTP方法
	Enabled bool          // 是否启用
}

// cacheResponseWriter 包装响应写入器以捕获响应数据
type cacheResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *cacheResponseWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// NewCacheMiddleware 创建缓存中间件
func NewCacheMiddleware() *CacheMiddleware {
	config := &CacheConfig{
		Enabled:    true,
		DefaultTTL: 10 * time.Minute,
		CachePaths: map[string]CachePathConfig{
			// 文章相关接口缓存配置
			"/api/articles": {
				TTL:     5 * time.Minute,
				Methods: []string{"GET"},
				Enabled: true,
			},
			"/api/articles/:id": {
				TTL:     30 * time.Minute,
				Methods: []string{"GET"},
				Enabled: true,
			},
			// 分类相关接口缓存配置
			"/api/categories": {
				TTL:     1 * time.Hour,
				Methods: []string{"GET"},
				Enabled: true,
			},
			"/api/categories/:id": {
				TTL:     1 * time.Hour,
				Methods: []string{"GET"},
				Enabled: true,
			},
			"/api/categories/:id/articles": {
				TTL:     10 * time.Minute,
				Methods: []string{"GET"},
				Enabled: true,
			},
			// 用户相关接口缓存配置
			"/api/users/:id/profile": {
				TTL:     30 * time.Minute,
				Methods: []string{"GET"},
				Enabled: true,
			},
			// 热门内容缓存配置
			"/api/articles/hot": {
				TTL:     6 * time.Hour,
				Methods: []string{"GET"},
				Enabled: true,
			},
			"/api/articles/recent": {
				TTL:     5 * time.Minute,
				Methods: []string{"GET"},
				Enabled: true,
			},
		},
	}

	return &CacheMiddleware{
		strategy: cache.NewCacheStrategy(),
		config:   config,
	}
}

// CacheHandler 缓存处理中间件
func (cm *CacheMiddleware) CacheHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查缓存是否启用
		if !cm.config.Enabled {
			c.Next()
			return
		}

		// 检查请求方法和路径是否需要缓存
		if !cm.shouldCache(c) {
			c.Next()
			return
		}

		// 生成缓存键
		cacheKey := cm.generateCacheKey(c)

		// 尝试从缓存获取数据
		cachedData, err := cm.strategy.GetArticleList(cacheKey)
		if err == nil && cachedData != nil {
			// 缓存命中，直接返回
			logger.Info("缓存命中",
				logger.String("key", cacheKey),
				logger.String("path", c.Request.URL.Path),
			)

			c.Header("X-Cache", "HIT")
			c.JSON(http.StatusOK, cachedData)
			c.Abort()
			return
		}

		// 缓存未命中，包装响应写入器
		writer := &cacheResponseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer

		// 继续处理请求
		c.Next()

		// 如果响应成功，缓存响应数据
		if c.Writer.Status() == http.StatusOK {
			cm.cacheResponse(cacheKey, writer.body.Bytes(), c)
		}
	}
}

// shouldCache 判断是否应该缓存
func (cm *CacheMiddleware) shouldCache(c *gin.Context) bool {
	path := c.FullPath()
	method := c.Request.Method

	// 检查路径配置
	pathConfig, exists := cm.config.CachePaths[path]
	if !exists {
		return false
	}

	if !pathConfig.Enabled {
		return false
	}

	// 检查HTTP方法
	for _, allowedMethod := range pathConfig.Methods {
		if method == allowedMethod {
			return true
		}
	}

	return false
}

// generateCacheKey 生成缓存键
func (cm *CacheMiddleware) generateCacheKey(c *gin.Context) string {
	// 基础路径
	path := c.FullPath()
	
	// 路径参数
	params := make(map[string]string)
	for _, param := range c.Params {
		params[param.Key] = param.Value
	}
	
	// 查询参数
	query := c.Request.URL.RawQuery
	
	// 用户ID（如果已认证）
	userID := ""
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(map[string]interface{}); ok {
			if id, exists := u["id"]; exists {
				userID = fmt.Sprintf("%v", id)
			}
		}
	}

	// 生成缓存键
	keyData := map[string]interface{}{
		"path":   path,
		"params": params,
		"query":  query,
		"user":   userID,
	}

	keyBytes, _ := json.Marshal(keyData)
	hash := md5.Sum(keyBytes)
	
	return fmt.Sprintf("api:cache:%x", hash)
}

// cacheResponse 缓存响应数据
func (cm *CacheMiddleware) cacheResponse(cacheKey string, responseData []byte, c *gin.Context) {
	// 解析响应数据
	var responseObj response.Response
	if err := json.Unmarshal(responseData, &responseObj); err != nil {
		logger.Error("解析响应数据失败",
			logger.String("key", cacheKey),
			logger.Err("error", err),
		)
		return
	}

	// 只缓存成功的响应
	if responseObj.Code != 0 {
		return
	}

	// 获取缓存时间
	ttl := cm.getCacheTTL(c.FullPath())

	// 缓存数据
	err := cm.strategy.SetArticleList(cacheKey, responseObj)
	if err != nil {
		logger.Error("缓存响应数据失败",
			logger.String("key", cacheKey),
			logger.Err("error", err),
		)
		return
	}

	logger.Info("缓存响应数据成功",
		logger.String("key", cacheKey),
		logger.String("path", c.Request.URL.Path),
		logger.String("ttl", ttl.String()),
	)

	// 设置缓存头
	c.Header("X-Cache", "MISS")
}

// getCacheTTL 获取缓存时间
func (cm *CacheMiddleware) getCacheTTL(path string) time.Duration {
	if pathConfig, exists := cm.config.CachePaths[path]; exists {
		return pathConfig.TTL
	}
	return cm.config.DefaultTTL
}

// InvalidateCache 缓存失效中间件
func (cm *CacheMiddleware) InvalidateCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理请求
		c.Next()

		// 如果是修改操作且成功，清除相关缓存
		if cm.shouldInvalidateCache(c) {
			cm.invalidateRelatedCache(c)
		}
	}
}

// shouldInvalidateCache 判断是否应该清除缓存
func (cm *CacheMiddleware) shouldInvalidateCache(c *gin.Context) bool {
	method := c.Request.Method
	status := c.Writer.Status()

	// 只在成功的修改操作后清除缓存
	if status < 200 || status >= 300 {
		return false
	}

	// 检查是否是修改操作
	modifyMethods := []string{"POST", "PUT", "PATCH", "DELETE"}
	for _, m := range modifyMethods {
		if method == m {
			return true
		}
	}

	return false
}

// invalidateRelatedCache 清除相关缓存
func (cm *CacheMiddleware) invalidateRelatedCache(c *gin.Context) {
	path := c.FullPath()
	method := c.Request.Method

	logger.Info("开始清除相关缓存",
		logger.String("path", path),
		logger.String("method", method),
	)

	// 根据路径清除相关缓存
	switch {
	case strings.Contains(path, "/articles"):
		cm.invalidateArticleCache(c)
	case strings.Contains(path, "/categories"):
		cm.invalidateCategoryCache(c)
	case strings.Contains(path, "/users"):
		cm.invalidateUserCache(c)
	}
}

// invalidateArticleCache 清除文章相关缓存
func (cm *CacheMiddleware) invalidateArticleCache(c *gin.Context) {
	// 获取文章ID
	if articleIDStr := c.Param("id"); articleIDStr != "" {
		if articleID, err := strconv.ParseUint(articleIDStr, 10, 32); err == nil {
			cm.strategy.InvalidateArticleCache(uint(articleID))
		}
	}

	// 清除文章列表缓存
	patterns := []string{
		"api:cache:*articles*",
		"articles:*",
	}

	for _, pattern := range patterns {
		if err := cm.strategy.DeletePattern(pattern); err != nil {
			logger.Error("清除文章缓存失败",
				logger.String("pattern", pattern),
				logger.Err("error", err),
			)
		}
	}
}

// invalidateCategoryCache 清除分类相关缓存
func (cm *CacheMiddleware) invalidateCategoryCache(c *gin.Context) {
	// 获取分类ID
	if categoryIDStr := c.Param("id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			cm.strategy.InvalidateCategoryCache(uint(categoryID))
		}
	}

	// 清除分类相关缓存
	patterns := []string{
		"api:cache:*categories*",
		"categories:*",
	}

	for _, pattern := range patterns {
		if err := cm.strategy.DeletePattern(pattern); err != nil {
			logger.Error("清除分类缓存失败",
				logger.String("pattern", pattern),
				logger.Err("error", err),
			)
		}
	}
}

// invalidateUserCache 清除用户相关缓存
func (cm *CacheMiddleware) invalidateUserCache(c *gin.Context) {
	// 获取用户ID
	if userIDStr := c.Param("id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			cm.strategy.InvalidateUserCache(uint(userID))
		}
	}

	// 清除用户相关缓存
	patterns := []string{
		"api:cache:*users*",
		"user:*",
	}

	for _, pattern := range patterns {
		if err := cm.strategy.DeletePattern(pattern); err != nil {
			logger.Error("清除用户缓存失败",
				logger.String("pattern", pattern),
				logger.Err("error", err),
			)
		}
	}
}

// GetCacheStats 获取缓存统计信息
func (cm *CacheMiddleware) GetCacheStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以实现缓存统计信息的获取
		stats := map[string]interface{}{
			"enabled":     cm.config.Enabled,
			"cache_paths": len(cm.config.CachePaths),
			"default_ttl": cm.config.DefaultTTL.String(),
		}

		response.Success(c, stats)
	}
}