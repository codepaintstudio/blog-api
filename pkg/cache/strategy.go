package cache

import (
	"fmt"
	"time"
)

// 缓存键常量定义
const (
	// 用户相关缓存键
	UserInfoKey     = "user:%d"           // 用户信息
	UserStatsKey    = "user:stats:%d"     // 用户统计
	UserTokenKey    = "user:token:%d"     // 用户令牌黑名单

	// 文章相关缓存键
	ArticleKey      = "article:%d"        // 文章详情
	ArticleListKey  = "articles:list:%s"  // 文章列表（按条件hash）
	HotArticlesKey  = "articles:hot"      // 热门文章
	RecentArticlesKey = "articles:recent" // 最新文章

	// 分类相关缓存键
	CategoryListKey = "categories:all"    // 所有分类
	CategoryKey     = "category:%d"       // 分类详情
	CategoryStatsKey = "category:stats:%d" // 分类统计

	// 统计相关缓存键
	ArticleLikesKey = "article:likes:%d"  // 文章点赞数
	ArticleViewsKey = "article:views:%d"  // 文章浏览量
	CommentCountKey = "article:comments:%d" // 文章评论数

	// API限流相关缓存键
	RateLimitKey    = "ratelimit:%s:%s"   // 限流键(IP/用户ID)
	ApiStatsKey     = "api:stats:%s:%s:%s" // API统计(endpoint:method:date)

	// 搜索相关缓存键
	SearchResultKey = "search:%s:%d:%d"   // 搜索结果(关键词:页码:大小)
	PopularSearchKey = "search:popular"   // 热门搜索关键词

	// 系统统计缓存键
	SystemStatsKey  = "system:stats:%s"   // 系统统计(日期)
	OnlineUsersKey  = "system:online"     // 在线用户数
)

// 缓存TTL配置
var (
	// 用户数据TTL
	UserInfoTTL     = 1 * time.Hour    // 用户信息1小时
	UserStatsTTL    = 30 * time.Minute // 用户统计30分钟
	UserTokenTTL    = 24 * time.Hour   // 用户令牌黑名单24小时

	// 文章数据TTL
	ArticleTTL      = 30 * time.Minute // 文章详情30分钟
	ArticleListTTL  = 10 * time.Minute // 文章列表10分钟
	HotArticlesTTL  = 6 * time.Hour    // 热门文章6小时
	RecentArticlesTTL = 5 * time.Minute // 最新文章5分钟

	// 分类数据TTL
	CategoryTTL     = 24 * time.Hour   // 分类详情24小时
	CategoryListTTL = 24 * time.Hour   // 分类列表24小时
	CategoryStatsTTL = 1 * time.Hour   // 分类统计1小时

	// 统计数据TTL
	ArticleLikesTTL = 5 * time.Minute  // 文章点赞数5分钟
	ArticleViewsTTL = 10 * time.Minute // 文章浏览量10分钟
	CommentCountTTL = 5 * time.Minute  // 评论数5分钟

	// API相关TTL
	ApiStatsTTL     = 1 * time.Hour    // API统计1小时
	
	// 搜索相关TTL
	SearchResultTTL = 15 * time.Minute // 搜索结果15分钟
	PopularSearchTTL = 1 * time.Hour   // 热门搜索1小时

	// 系统统计TTL
	SystemStatsTTL  = 1 * time.Hour    // 系统统计1小时
	OnlineUsersTTL  = 1 * time.Minute  // 在线用户数1分钟
)

// CacheStrategy 缓存策略
type CacheStrategy struct {
	service CacheService
	keyGen  *CacheKey
}

// NewCacheStrategy 创建缓存策略
func NewCacheStrategy() *CacheStrategy {
	return &CacheStrategy{
		service: NewCacheService(),
		keyGen:  NewCacheKey(),
	}
}

// GetUserInfo 获取用户信息缓存
func (cs *CacheStrategy) GetUserInfo(userID uint) (interface{}, error) {
	key := fmt.Sprintf(UserInfoKey, userID)
	var result interface{}
	err := cs.service.Get(key, &result)
	return result, err
}

// SetUserInfo 设置用户信息缓存
func (cs *CacheStrategy) SetUserInfo(userID uint, userInfo interface{}) error {
	key := fmt.Sprintf(UserInfoKey, userID)
	return cs.service.Set(key, userInfo, UserInfoTTL)
}

// GetArticle 获取文章缓存
func (cs *CacheStrategy) GetArticle(articleID uint) (interface{}, error) {
	key := fmt.Sprintf(ArticleKey, articleID)
	var result interface{}
	err := cs.service.Get(key, &result)
	return result, err
}

// SetArticle 设置文章缓存
func (cs *CacheStrategy) SetArticle(articleID uint, article interface{}) error {
	key := fmt.Sprintf(ArticleKey, articleID)
	return cs.service.Set(key, article, ArticleTTL)
}

// GetArticleList 获取文章列表缓存
func (cs *CacheStrategy) GetArticleList(cacheKey string) (interface{}, error) {
	key := fmt.Sprintf(ArticleListKey, cacheKey)
	var result interface{}
	err := cs.service.Get(key, &result)
	return result, err
}

// SetArticleList 设置文章列表缓存
func (cs *CacheStrategy) SetArticleList(cacheKey string, articles interface{}) error {
	key := fmt.Sprintf(ArticleListKey, cacheKey)
	return cs.service.Set(key, articles, ArticleListTTL)
}

// GetHotArticles 获取热门文章缓存
func (cs *CacheStrategy) GetHotArticles() (interface{}, error) {
	var result interface{}
	err := cs.service.Get(HotArticlesKey, &result)
	return result, err
}

// SetHotArticles 设置热门文章缓存
func (cs *CacheStrategy) SetHotArticles(articles interface{}) error {
	return cs.service.Set(HotArticlesKey, articles, HotArticlesTTL)
}

// GetCategoryList 获取分类列表缓存
func (cs *CacheStrategy) GetCategoryList() (interface{}, error) {
	var result interface{}
	err := cs.service.Get(CategoryListKey, &result)
	return result, err
}

// SetCategoryList 设置分类列表缓存
func (cs *CacheStrategy) SetCategoryList(categories interface{}) error {
	return cs.service.Set(CategoryListKey, categories, CategoryListTTL)
}

// IncrArticleViews 增加文章浏览量
func (cs *CacheStrategy) IncrArticleViews(articleID uint) (int64, error) {
	key := fmt.Sprintf(ArticleViewsKey, articleID)
	count, err := cs.service.Incr(key)
	if err != nil {
		return 0, err
	}
	// 设置过期时间
	cs.service.Expire(key, ArticleViewsTTL)
	return count, nil
}

// GetArticleViews 获取文章浏览量
func (cs *CacheStrategy) GetArticleViews(articleID uint) (int64, error) {
	key := fmt.Sprintf(ArticleViewsKey, articleID)
	var count int64
	err := cs.service.Get(key, &count)
	return count, err
}

// InvalidateUserCache 清除用户相关缓存
func (cs *CacheStrategy) InvalidateUserCache(userID uint) error {
	patterns := []string{
		fmt.Sprintf("user:%d*", userID),
	}
	
	for _, pattern := range patterns {
		if err := cs.service.DeletePattern(pattern); err != nil {
			return err
		}
	}
	return nil
}

// InvalidateArticleCache 清除文章相关缓存
func (cs *CacheStrategy) InvalidateArticleCache(articleID uint) error {
	patterns := []string{
		fmt.Sprintf("article:%d*", articleID),
		"articles:*", // 清除所有文章列表缓存
		HotArticlesKey,
		RecentArticlesKey,
	}
	
	for _, pattern := range patterns {
		if err := cs.service.DeletePattern(pattern); err != nil {
			return err
		}
	}
	return nil
}

// InvalidateCategoryCache 清除分类相关缓存
func (cs *CacheStrategy) InvalidateCategoryCache(categoryID uint) error {
	patterns := []string{
		fmt.Sprintf("category:%d*", categoryID),
		CategoryListKey,
		"articles:*", // 分类变更影响文章列表
	}
	
	for _, pattern := range patterns {
		if err := cs.service.DeletePattern(pattern); err != nil {
			return err
		}
	}
	return nil
}

// GetSearchResults 获取搜索结果缓存
func (cs *CacheStrategy) GetSearchResults(keyword string, page, size int) (interface{}, error) {
	key := fmt.Sprintf(SearchResultKey, keyword, page, size)
	var result interface{}
	err := cs.service.Get(key, &result)
	return result, err
}

// SetSearchResults 设置搜索结果缓存
func (cs *CacheStrategy) SetSearchResults(keyword string, page, size int, results interface{}) error {
	key := fmt.Sprintf(SearchResultKey, keyword, page, size)
	return cs.service.Set(key, results, SearchResultTTL)
}

// SetRateLimit 设置限流缓存
func (cs *CacheStrategy) SetRateLimit(identifier, action string, count int64, ttl time.Duration) error {
	key := fmt.Sprintf(RateLimitKey, identifier, action)
	return cs.service.Set(key, count, ttl)
}

// GetRateLimit 获取限流计数
func (cs *CacheStrategy) GetRateLimit(identifier, action string) (int64, error) {
	key := fmt.Sprintf(RateLimitKey, identifier, action)
	var count int64
	err := cs.service.Get(key, &count)
	return count, err
}

// IncrRateLimit 增加限流计数
func (cs *CacheStrategy) IncrRateLimit(identifier, action string, ttl time.Duration) (int64, error) {
	key := fmt.Sprintf(RateLimitKey, identifier, action)
	count, err := cs.service.Incr(key)
	if err != nil {
		return 0, err
	}
	
	// 如果是第一次设置，设置过期时间
	if count == 1 {
		cs.service.Expire(key, ttl)
	}
	
	return count, nil
}

// DeletePattern 删除匹配模式的缓存
func (cs *CacheStrategy) DeletePattern(pattern string) error {
	return cs.service.DeletePattern(pattern)
}