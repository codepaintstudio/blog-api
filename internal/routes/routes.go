package routes

import (
	"time"
	
	"blog-api/pkg/config"
	"blog-api/pkg/response"
	"blog-api/internal/middlewares"
	"blog-api/internal/controllers"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	
	_ "blog-api/docs" // 导入生成的swagger文档
)

// SetupRoutes 设置路由
func SetupRoutes() *gin.Engine {
	cfg := config.GetConfig()
	
	// 创建Gin实例
	r := gin.New()
	
	// 添加中间件
	setupMiddlewares(r, cfg)
	
	// 设置API路由
	setupAPIRoutes(r)
	
	// 设置Swagger文档路由
	setupSwagger(r)
	
	return r
}

// setupMiddlewares 设置中间件
func setupMiddlewares(r *gin.Engine, cfg *config.Config) {
	// 自定义Recovery中间件
	r.Use(middlewares.RecoveryMiddleware())
	
	// 自定义日志中间件
	r.Use(middlewares.LoggerMiddleware())
	
	// 限流中间件
	middlewares.InitRateLimiter(&cfg.RateLimit)
	r.Use(middlewares.RateLimitMiddleware())
	
	// CORS中间件
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		ExposeHeaders:    cfg.CORS.ExposeHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           time.Duration(cfg.CORS.MaxAge) * time.Hour,
	}
	r.Use(cors.New(corsConfig))
}

// setupAPIRoutes 设置API路由
func setupAPIRoutes(r *gin.Engine) {
	// 初始化控制器
	authController := controllers.NewAuthController()
	userController := controllers.NewUserController()
	articleController := controllers.NewArticleController()
	categoryController := controllers.NewCategoryController()
	
	// API版本分组
	v1 := r.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", healthCheck)
		
		// 认证路由（无需登录）
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authController.Register)
			authGroup.POST("/login", authController.Login)
			authGroup.POST("/refresh", authController.RefreshToken)
		}
		
		// 用户路由（需要登录）
		userGroup := v1.Group("/users")
		userGroup.Use(middlewares.AuthMiddleware())
		{
			userGroup.GET("/profile", userController.GetProfile)
			userGroup.PUT("/profile", userController.UpdateProfile)
			userGroup.PUT("/password", userController.ChangePassword)
		}
		
		// 文章路由
		articleGroup := v1.Group("/articles")
		{
			// 公开路由（无需登录）
			articleGroup.GET("", articleController.List)           // 获取文章列表
			articleGroup.GET("/search", articleController.Search)   // 搜索文章
			articleGroup.GET("/:id", articleController.GetByID)     // 获取文章详情
			
			// 需要登录的路由
			authRequired := articleGroup.Group("")
			authRequired.Use(middlewares.AuthMiddleware())
			{
				authRequired.POST("", articleController.Create)        // 创建文章
				authRequired.PUT("/:id", articleController.Update)     // 更新文章
				authRequired.DELETE("/:id", articleController.Delete)  // 删除文章
			}
		}
		
		// 分类路由
		categoryGroup := v1.Group("/categories")
		{
			// 公开路由（无需登录）
			categoryGroup.GET("", categoryController.List)         // 获取所有分类列表
			categoryGroup.GET("/active", categoryController.ListActive)  // 获取活跃分类列表
			categoryGroup.GET("/:id", categoryController.GetByID)   // 获取分类详情
			
			// 需要登录的路由
			authRequired := categoryGroup.Group("")
			authRequired.Use(middlewares.AuthMiddleware())
			{
				authRequired.POST("", categoryController.Create)       // 创建分类
				authRequired.PUT("/:id", categoryController.Update)    // 更新分类
				authRequired.DELETE("/:id", categoryController.Delete) // 删除分类
			}
		}
	}
}

// setupSwagger 设置Swagger文档路由
func setupSwagger(r *gin.Engine) {
	// Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// 重定向根路径到Swagger文档
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})
}

// healthCheck 健康检查接口
// @Summary 健康检查
// @Description 检查服务是否正常运行
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /health [get]
func healthCheck(c *gin.Context) {
	response.SuccessWithMessage(c, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
	}, "服务运行正常")
}