package middleware

import (
	"net"
	"net/http"
	"server/pkg/response"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter IP限流器
type IPRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit // 每秒允许的请求数
	burst    int        // 突发请求数
	
	// IP封禁相关
	bannedIPs     map[string]time.Time // 被封禁的IP及封禁到期时间
	banMu         sync.RWMutex
	violationMap  map[string]int       // IP违规次数统计
	violationMu   sync.RWMutex
	banDuration   time.Duration        // 封禁持续时间
	maxViolations int                  // 最大违规次数
}

// NewIPRateLimiter 创建新的IP限流器
// r: 每秒允许的请求数 (例如: 10 表示每秒10个请求)
// b: 突发请求数 (例如: 20 表示突发时允许20个请求)
// banDuration: IP封禁时长
// maxViolations: 触发封禁的违规次数
func NewIPRateLimiter(r rate.Limit, b int, banDuration time.Duration, maxViolations int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		limiters:      make(map[string]*rate.Limiter),
		rate:          r,
		burst:         b,
		bannedIPs:     make(map[string]time.Time),
		violationMap:  make(map[string]int),
		banDuration:   banDuration,
		maxViolations: maxViolations,
	}
	
	// 启动清理协程，定期清理过期的限流器和封禁记录
	go limiter.cleanupRoutine()
	
	return limiter
}

// getLimiter 获取IP对应的限流器
func (ipl *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	ipl.mu.Lock()
	defer ipl.mu.Unlock()
	
	limiter, exists := ipl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(ipl.rate, ipl.burst)
		ipl.limiters[ip] = limiter
	}
	
	return limiter
}

// isIPBanned 检查IP是否被封禁
func (ipl *IPRateLimiter) isIPBanned(ip string) bool {
	ipl.banMu.RLock()
	defer ipl.banMu.RUnlock()
	
	banTime, exists := ipl.bannedIPs[ip]
	if !exists {
		return false
	}
	
	// 检查封禁是否过期
	if time.Now().After(banTime) {
		// 封禁过期，删除记录
		go func() {
			ipl.banMu.Lock()
			delete(ipl.bannedIPs, ip)
			ipl.banMu.Unlock()
		}()
		return false
	}
	
	return true
}

// recordViolation 记录IP违规
func (ipl *IPRateLimiter) recordViolation(ip string) {
	ipl.violationMu.Lock()
	defer ipl.violationMu.Unlock()
	
	ipl.violationMap[ip]++
	
	// 检查是否达到封禁阈值
	if ipl.violationMap[ip] >= ipl.maxViolations {
		ipl.banMu.Lock()
		ipl.bannedIPs[ip] = time.Now().Add(ipl.banDuration)
		ipl.banMu.Unlock()
		
		// 重置违规计数
		ipl.violationMap[ip] = 0
	}
}

// getClientIP 获取客户端真实IP
func getClientIP(c *gin.Context) string {
	// 检查代理头
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		// X-Forwarded-For 可能包含多个IP，取第一个
		if i := len(ip); i > 0 {
			if idx := 0; idx < i {
				for idx < i && ip[idx] == ' ' {
					idx++
				}
				for end := idx; end < i && ip[end] != ',' && ip[end] != ' '; end++ {
					if end > idx {
						return ip[idx:end]
					}
				}
			}
		}
		return ip
	}
	
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	
	// 获取远程地址
	if ip, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
		return ip
	}
	
	return c.Request.RemoteAddr
}

// RateLimitMiddleware 限流中间件
func (ipl *IPRateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)
		
		// 检查IP是否被封禁
		if ipl.isIPBanned(ip) {
			c.JSON(http.StatusTooManyRequests, response.Error(response.StatusTooManyRequests, "IP temporarily banned due to too many violations"))
			c.Abort()
			return
		}
		
		// 获取该IP的限流器
		limiter := ipl.getLimiter(ip)
		
		// 检查是否允许请求
		if !limiter.Allow() {
			// 记录违规
			ipl.recordViolation(ip)
			
			c.JSON(http.StatusTooManyRequests, response.Error(response.StatusTooManyRequests, "Rate limit exceeded"))
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// cleanupRoutine 清理过期记录
func (ipl *IPRateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(10 * time.Minute) // 每10分钟清理一次
	defer ticker.Stop()
	
	for range ticker.C {
		now := time.Now()
		
		// 清理过期的封禁记录
		ipl.banMu.Lock()
		for ip, banTime := range ipl.bannedIPs {
			if now.After(banTime) {
				delete(ipl.bannedIPs, ip)
			}
		}
		ipl.banMu.Unlock()
		
		// 清理长时间未使用的限流器 (1小时未使用)
		ipl.mu.Lock()
		for ip, limiter := range ipl.limiters {
			// 简单检查：如果限流器的令牌桶是满的，说明很久没有请求了
			if limiter.Tokens() == float64(ipl.burst) {
				delete(ipl.limiters, ip)
			}
		}
		ipl.mu.Unlock()
		
		// 清理违规记录 (1小时后重置)
		ipl.violationMu.Lock()
		// 简单实现：定期重置所有违规计数
		// 在生产环境中，可以为每个IP记录时间戳，更精确地清理
		if len(ipl.violationMap) > 1000 { // 当记录过多时清理
			ipl.violationMap = make(map[string]int)
		}
		ipl.violationMu.Unlock()
	}
}

// GetStats 获取限流统计信息
func (ipl *IPRateLimiter) GetStats() map[string]interface{} {
	ipl.mu.RLock()
	activeIPs := len(ipl.limiters)
	ipl.mu.RUnlock()
	
	ipl.banMu.RLock()
	bannedCount := len(ipl.bannedIPs)
	ipl.banMu.RUnlock()
	
	return map[string]interface{}{
		"active_ips":   activeIPs,
		"banned_count": bannedCount,
		"rate_limit":   float64(ipl.rate),
		"burst_size":   ipl.burst,
	}
}