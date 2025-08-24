package routes

import (
	"server/internal/controllers"
	"server/pkg/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// 创建IP限流器
	// 参数：每秒20个请求，突发30个请求，封禁30分钟，5次违规后封禁
	ipLimiter := middleware.NewIPRateLimiter(
		rate.Every(50*time.Millisecond), // 每50ms一个请求 = 每秒20个请求
		30,                                // 突发请求数
		10*time.Minute,                    // IP封禁时长
		5,                                 // 最大违规次数
	)

	// 使用中间件
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(ipLimiter.RateLimitMiddleware()) // 添加限流中间件

	// 创建控制器实例
	userController := controllers.NewUserController(db)
	articleController := controllers.NewArticleController(db)

	// API 路由组
	api := router.Group("/api")

	// 系统监控接口（可选，用于查看限流状态）
	api.GET("/system/rate-limit-stats", func(c *gin.Context) {
		stats := ipLimiter.GetStats()
		c.JSON(200, gin.H{
			"code":    200,
			"message": "Rate limit stats",
			"data":    stats,
		})
	})

	// 用户相关路由
	user := api.Group("/user")
	{
		user.POST("/register", userController.Register)
		user.POST("/login", userController.Login)
	}

	// 用户信息查询路由
	users := api.Group("/users")
	{
		users.GET("/:id", userController.GetUserById)          // 获取用户基本信息
		users.GET("/:id/detail", userController.GetUserDetail) // 获取用户详情（含统计）
	}

	// 帖子相关路由
	article := api.Group("/articles")
	{
		// 公开路由
		article.GET("", articleController.List)           // 帖子列表（支持搜索、排序、过滤）
		article.GET("/:id", articleController.GetById)    // 帖子详情
		article.GET("/stats", articleController.GetStats) // 文章统计信息

		// 需要登录的路由
		auth := article.Group("", middleware.AuthMiddleware())
		{
			auth.POST("", articleController.Create)       // 创建帖子
			auth.PUT("/:id", articleController.Update)    // 更新帖子
			auth.DELETE("/:id", articleController.Delete) // 删除帖子
		}
	}
}
