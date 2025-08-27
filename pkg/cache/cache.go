package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"blog-api/pkg/database"
	"blog-api/pkg/logger"
	"github.com/go-redis/redis/v8"
)

// CacheService 缓存服务接口
type CacheService interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string, dest interface{}) error
	Delete(key string) error
	DeletePattern(pattern string) error
	Exists(key string) bool
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	Incr(key string) (int64, error)
	IncrBy(key string, value int64) (int64, error)
	Expire(key string, expiration time.Duration) error
}

// cacheService 缓存服务实现
type cacheService struct {
	client *redis.Client
	ctx    context.Context
}

// NewCacheService 创建缓存服务
func NewCacheService() CacheService {
	return &cacheService{
		client: database.GetRedis(),
		ctx:    context.Background(),
	}
}

// Set 设置缓存
func (c *cacheService) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化缓存数据失败: %w", err)
	}

	err = c.client.Set(c.ctx, key, data, expiration).Err()
	if err != nil {
		logger.Error("设置缓存失败",
			logger.String("key", key),
			logger.Err("error", err),
		)
		return err
	}

	logger.Debug("缓存设置成功",
		logger.String("key", key),
		logger.String("expiration", expiration.String()),
	)
	return nil
}

// Get 获取缓存
func (c *cacheService) Get(key string, dest interface{}) error {
	data, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("缓存未找到: %s", key)
	}
	if err != nil {
		logger.Error("获取缓存失败",
			logger.String("key", key),
			logger.Err("error", err),
		)
		return err
	}

	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		return fmt.Errorf("反序列化缓存数据失败: %w", err)
	}

	logger.Debug("缓存获取成功", logger.String("key", key))
	return nil
}

// Delete 删除缓存
func (c *cacheService) Delete(key string) error {
	err := c.client.Del(c.ctx, key).Err()
	if err != nil {
		logger.Error("删除缓存失败",
			logger.String("key", key),
			logger.Err("error", err),
		)
		return err
	}

	logger.Debug("缓存删除成功", logger.String("key", key))
	return nil
}

// DeletePattern 删除匹配模式的缓存
func (c *cacheService) DeletePattern(pattern string) error {
	keys, err := c.client.Keys(c.ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	err = c.client.Del(c.ctx, keys...).Err()
	if err != nil {
		logger.Error("批量删除缓存失败",
			logger.String("pattern", pattern),
			logger.Int("count", len(keys)),
			logger.Err("error", err),
		)
		return err
	}

	logger.Debug("批量删除缓存成功",
		logger.String("pattern", pattern),
		logger.Int("count", len(keys)),
	)
	return nil
}

// Exists 检查缓存是否存在
func (c *cacheService) Exists(key string) bool {
	count, err := c.client.Exists(c.ctx, key).Result()
	if err != nil {
		logger.Error("检查缓存存在性失败",
			logger.String("key", key),
			logger.Err("error", err),
		)
		return false
	}
	return count > 0
}

// SetNX 仅在键不存在时设置缓存 (分布式锁)
func (c *cacheService) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("序列化缓存数据失败: %w", err)
	}

	result, err := c.client.SetNX(c.ctx, key, data, expiration).Result()
	if err != nil {
		logger.Error("SetNX操作失败",
			logger.String("key", key),
			logger.Err("error", err),
		)
		return false, err
	}

	return result, nil
}

// Incr 递增计数器
func (c *cacheService) Incr(key string) (int64, error) {
	result, err := c.client.Incr(c.ctx, key).Result()
	if err != nil {
		logger.Error("递增计数器失败",
			logger.String("key", key),
			logger.Err("error", err),
		)
		return 0, err
	}
	return result, nil
}

// IncrBy 按指定值递增计数器
func (c *cacheService) IncrBy(key string, value int64) (int64, error) {
	result, err := c.client.IncrBy(c.ctx, key, value).Result()
	if err != nil {
		logger.Error("递增计数器失败",
			logger.String("key", key),
			logger.String("value", fmt.Sprintf("%d", value)),
			logger.Err("error", err),
		)
		return 0, err
	}
	return result, nil
}

// Expire 设置键的过期时间
func (c *cacheService) Expire(key string, expiration time.Duration) error {
	err := c.client.Expire(c.ctx, key, expiration).Err()
	if err != nil {
		logger.Error("设置键过期时间失败",
			logger.String("key", key),
			logger.String("expiration", expiration.String()),
			logger.Err("error", err),
		)
		return err
	}
	return nil
}

// CacheKey 构建缓存键的工具函数
type CacheKey struct{}

// NewCacheKey 创建缓存键构建器
func NewCacheKey() *CacheKey {
	return &CacheKey{}
}

// UserProfile 用户档案缓存键
func (ck *CacheKey) UserProfile(userID uint) string {
	return fmt.Sprintf("user:profile:%d", userID)
}

// Article 文章缓存键
func (ck *CacheKey) Article(articleID uint) string {
	return fmt.Sprintf("article:%d", articleID)
}

// ArticleList 文章列表缓存键
func (ck *CacheKey) ArticleList(page, size int, categoryID uint, status, keyword string) string {
	return fmt.Sprintf("articles:list:%d:%d:%d:%s:%s", page, size, categoryID, status, keyword)
}

// Category 分类缓存键
func (ck *CacheKey) Category(categoryID uint) string {
	return fmt.Sprintf("category:%d", categoryID)
}

// CategoryList 分类列表缓存键
func (ck *CacheKey) CategoryList() string {
	return "categories:list"
}

// UserFavorites 用户收藏缓存键
func (ck *CacheKey) UserFavorites(userID uint, folderID uint) string {
	return fmt.Sprintf("user:%d:favorites:%d", userID, folderID)
}

// PopularArticles 热门文章缓存键
func (ck *CacheKey) PopularArticles() string {
	return "articles:popular"
}

// ViewCount 浏览量缓存键
func (ck *CacheKey) ViewCount(articleID uint) string {
	return fmt.Sprintf("article:%d:views", articleID)
}

// ApiStats API统计缓存键
func (ck *CacheKey) ApiStats(endpoint, method string, date string) string {
	return fmt.Sprintf("api:stats:%s:%s:%s", endpoint, method, date)
}