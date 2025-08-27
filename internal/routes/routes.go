package routes

import (
	"time"
	
	"blog-api/pkg/config"
	"blog-api/pkg/response"
	
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
	// Recovery中间件
	r.Use(gin.Recovery())
	
	// 日志中间件
	r.Use(gin.Logger())
	
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
	// API版本分组
	v1 := r.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", healthCheck)
		
		// 认证路由
		authGroup := v1.Group("/auth")
		{
			// 这里后续添加认证相关路由
			authGroup.GET("/test", func(c *gin.Context) {
				response.Success(c, gin.H{"message": "auth test"})
			})
		}
		
		// 用户路由
		userGroup := v1.Group("/users")
		{
			// 这里后续添加用户相关路由
			userGroup.GET("/test", func(c *gin.Context) {
				response.Success(c, gin.H{"message": "user test"})
			})
		}
		
		// 文章路由
		articleGroup := v1.Group("/articles")
		{
			// 这里后续添加文章相关路由
			articleGroup.GET("/test", func(c *gin.Context) {
				response.Success(c, gin.H{"message": "article test"})
			})
		}
		
		// 分类路由
		categoryGroup := v1.Group("/categories")
		{
			// 这里后续添加分类相关路由
			categoryGroup.GET("/test", func(c *gin.Context) {
				response.Success(c, gin.H{"message": "category test"})
			})
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