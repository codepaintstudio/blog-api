package middlewares

import (
	"sync"
	"time"

	"blog-api/pkg/config"
	"blog-api/pkg/logger"
	"blog-api/pkg/response"
	
	"github.com/gin-gonic/gin"
)

// RateLimiter 限流器结构
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // 每分钟允许的请求数
	burst    int           // 突发请求数
	cleanup  time.Duration // 清理间隔
}

// Visitor 访客信息
type Visitor struct {
	limiter  *TokenBucket
	lastSeen time.Time
}

// TokenBucket 令牌桶
type TokenBucket struct {
	capacity int       // 桶容量
	tokens   int       // 当前令牌数
	rate     int       // 添加令牌的速率（每分钟）
	lastFill time.Time // 上次填充时间
	mu       sync.Mutex
}

// NewTokenBucket 创建新的令牌桶
func NewTokenBucket(capacity, rate int) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		tokens:   capacity,
		rate:     rate,
		lastFill: time.Now(),
	}
}

// Allow 检查是否允许请求
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(tb.lastFill)
	
	// 计算需要添加的令牌数
	tokensToAdd := int(elapsed.Minutes() * float64(tb.rate))
	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastFill = now
	}
	
	// 检查是否有可用令牌
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	
	return false
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		burst:    burst,
		cleanup:  5 * time.Minute,
	}
	
	// 启动清理goroutine
	go rl.cleanupVisitors()
	
	return rl
}

// getVisitor 获取或创建访客
func (rl *RateLimiter) getVisitor(ip string) *Visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	v, exists := rl.visitors[ip]
	if !exists {
		v = &Visitor{
			limiter:  NewTokenBucket(rl.burst, rl.rate),
			lastSeen: time.Now(),
		}
		rl.visitors[ip] = v
	}
	
	v.lastSeen = time.Now()
	return v
}

// cleanupVisitors 清理过期的访客
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(rl.cleanup)
		
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.cleanup {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow 检查IP是否允许请求
func (rl *RateLimiter) Allow(ip string) bool {
	visitor := rl.getVisitor(ip)
	return visitor.limiter.Allow()
}

var globalRateLimiter *RateLimiter

// InitRateLimiter 初始化全局限流器
func InitRateLimiter(cfg *config.RateLimitConfig) {
	if cfg.Enabled {
		globalRateLimiter = NewRateLimiter(cfg.RequestsPerMinute, cfg.Burst)
		logger.Info("限流器已启用",
			logger.Int("requests_per_minute", cfg.RequestsPerMinute),
			logger.Int("burst", cfg.Burst),
		)
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if globalRateLimiter == nil {
			c.Next()
			return
		}
		
		clientIP := c.ClientIP()
		
		if !globalRateLimiter.Allow(clientIP) {
			logger.Warn("请求被限流",
				logger.String("ip", clientIP),
				logger.String("path", c.Request.URL.Path),
				logger.String("method", c.Request.Method),
			)
			
			response.ErrorWithMessage(c, response.ERROR, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		
		c.Next()
	}
}