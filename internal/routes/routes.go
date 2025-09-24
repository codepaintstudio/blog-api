package routes

import (
	"time"

	"blog-api/internal/controllers"
	"blog-api/internal/middlewares"
	"blog-api/pkg/config"
	"blog-api/pkg/response"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(
	authController *controllers.AuthController,
	userController *controllers.UserController,
	articleController *controllers.ArticleController,
	categoryController *controllers.CategoryController,
	fileController *controllers.FileController,
	commentController *controllers.CommentController,
	likeController *controllers.LikeController,
	favoriteController *controllers.FavoriteController,
	adminController *controllers.AdminController,
) *gin.Engine {
	cfg := config.GetConfig()

	// 创建Gin实例
	r := gin.New()

	// 添加中间件
	setupMiddlewares(r, cfg)

	// 设置API路由
	setupAPIRoutes(r, authController, userController, articleController, categoryController, fileController, commentController, likeController, favoriteController, adminController)

	return r
}

// setupMiddlewares 设置中间件
func setupMiddlewares(r *gin.Engine, cfg *config.Config) {
	// 安全头中间件 (最先执行)
	r.Use(middlewares.SecurityHeaders())

	// 自定义Recovery中间件
	r.Use(middlewares.RecoveryMiddleware())

	// 自定义日志中间件
	r.Use(middlewares.LoggerMiddleware())

	// 限流中间件
	middlewares.InitRateLimiter(&cfg.RateLimit)
	r.Use(middlewares.RateLimitMiddleware())

	// 缓存中间件
	cacheMiddleware := middlewares.NewCacheMiddleware()
	r.Use(cacheMiddleware.CacheHandler())
	r.Use(cacheMiddleware.InvalidateCache())

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
func setupAPIRoutes(
	r *gin.Engine,
	authController *controllers.AuthController,
	userController *controllers.UserController,
	articleController *controllers.ArticleController,
	categoryController *controllers.CategoryController,
	fileController *controllers.FileController,
	commentController *controllers.CommentController,
	likeController *controllers.LikeController,
	favoriteController *controllers.FavoriteController,
	adminController *controllers.AdminController,
) {
	// 静态文件服务
	r.Static("/uploads", "./uploads")

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
			articleGroup.GET("", articleController.List)          // 获取文章列表
			articleGroup.GET("/search", articleController.Search) // 搜索文章
			articleGroup.GET("/:id", articleController.GetByID)   // 获取文章详情

			// 需要登录的路由
			authRequired := articleGroup.Group("")
			authRequired.Use(middlewares.AuthMiddleware())
			{
				authRequired.POST("", articleController.Create)       // 创建文章
				authRequired.PUT("/:id", articleController.Update)    // 更新文章
				authRequired.DELETE("/:id", articleController.Delete) // 删除文章
			}
		}

		// 分类路由
		categoryGroup := v1.Group("/categories")
		{
			// 公开路由（无需登录）
			categoryGroup.GET("", categoryController.List)              // 获取所有分类列表
			categoryGroup.GET("/active", categoryController.ListActive) // 获取活跃分类列表
			categoryGroup.GET("/:id", categoryController.GetByID)       // 获取分类详情

			// 管理员专有路由（需要管理员权限）
			adminRequired := categoryGroup.Group("")
			adminRequired.Use(middlewares.AdminAuthMiddleware())
			{
				adminRequired.POST("", categoryController.Create)       // 创建分类
				adminRequired.PUT("/:id", categoryController.Update)    // 更新分类
				adminRequired.DELETE("/:id", categoryController.Delete) // 删除分类
			}
		}

		// 文件管理路由
		fileGroup := v1.Group("/files")
		fileGroup.Use(middlewares.AuthMiddleware())
		{
			fileGroup.POST("/upload", fileController.Upload)    // 文件上传
			fileGroup.GET("", fileController.GetUserFiles)      // 获取用户文件列表
			fileGroup.GET("/:id", fileController.GetFile)       // 获取文件信息
			fileGroup.DELETE("/:id", fileController.DeleteFile) // 删除文件
		}

		// 评论路由
		commentGroup := v1.Group("/comments")
		{
			// 公开路由（无需登录）
			commentGroup.GET("/:id", commentController.GetComment)                // 获取评论详情
			commentGroup.GET("/:id/replies", commentController.GetCommentReplies) // 获取评论回复列表

			// 需要登录的路由
			authRequired := commentGroup.Group("")
			authRequired.Use(middlewares.AuthMiddleware())
			{
				authRequired.POST("", commentController.CreateComment)       // 创建评论
				authRequired.PUT("/:id", commentController.UpdateComment)    // 更新评论
				authRequired.DELETE("/:id", commentController.DeleteComment) // 删除评论
				authRequired.GET("/mine", commentController.GetUserComments) // 获取用户评论列表

				// 管理员功能（需要后续添加管理员权限中间件）
				authRequired.POST("/:id/approve", commentController.ApproveComment) // 审核通过评论
				authRequired.POST("/:id/reject", commentController.RejectComment)   // 拒绝评论
			}
		}

		// 点赞路由
		likeGroup := v1.Group("/likes")
		likeGroup.Use(middlewares.AuthMiddleware())
		{
			likeGroup.GET("/articles", likeController.GetUserArticleLikes) // 获取用户点赞的文章列表
			likeGroup.GET("/comments", likeController.GetUserCommentLikes) // 获取用户点赞的评论列表
		}

		// 收藏路由
		favoriteGroup := v1.Group("/favorites")
		favoriteGroup.Use(middlewares.AuthMiddleware())
		{
			// 收藏夹管理
			favoriteGroup.POST("/folders", favoriteController.CreateFolder)                  // 创建收藏夹
			favoriteGroup.GET("/folders", favoriteController.GetUserFolders)                 // 获取用户收藏夹列表
			favoriteGroup.PUT("/folders/:id", favoriteController.UpdateFolder)               // 更新收藏夹
			favoriteGroup.DELETE("/folders/:id", favoriteController.DeleteFolder)            // 删除收藏夹
			favoriteGroup.GET("/folders/:id/articles", favoriteController.GetFolderArticles) // 获取收藏夹文章列表

			// 文章收藏管理
			favoriteGroup.POST("", favoriteController.AddToFavorite)                    // 添加文章到收藏夹
			favoriteGroup.DELETE("/:article_id", favoriteController.RemoveFromFavorite) // 从收藏夹移除文章
			favoriteGroup.GET("/articles", favoriteController.GetUserFavorites)         // 获取用户所有收藏文章
			favoriteGroup.GET("/:id/status", favoriteController.CheckFavoriteStatus)    // 检查文章收藏状态
		}
	}

	// 在文章路由中添加评论和点赞相关路由
	v1.GET("/articles/:id/comments", commentController.GetArticleComments)                        // 获取文章评论列表
	v1.POST("/articles/:id/like", middlewares.AuthMiddleware(), likeController.ToggleArticleLike) // 切换文章点赞状态
	v1.GET("/articles/:id/likers", likeController.GetArticleLikers)                               // 获取文章点赞用户列表

	// 在评论路由中添加点赞相关路由
	v1.POST("/comments/:id/like", middlewares.AuthMiddleware(), likeController.ToggleCommentLike) // 切换评论点赞状态

	// 管理员路由
	adminGroup := v1.Group("/admin")
	adminGroup.Use(middlewares.AdminAuthMiddleware()) // 使用管理员认证中间件
	{
		// 用户管理
		adminGroup.GET("/users", adminController.GetUserList)                 // 获取用户列表
		adminGroup.GET("/users/:id", adminController.GetUserDetails)          // 获取用户详情
		adminGroup.PUT("/users/:id/status", adminController.UpdateUserStatus) // 更新用户状态
		adminGroup.PUT("/users/:id/role", adminController.UpdateUserRole)     // 更新用户角色
		adminGroup.DELETE("/users/:id", adminController.DeleteUser)           // 删除用户

		// 系统统计
		adminGroup.GET("/stats/system", adminController.GetSystemStats)   // 获取系统统计
		adminGroup.GET("/stats/users", adminController.GetUserStats)      // 获取用户统计
		adminGroup.GET("/stats/content", adminController.GetContentStats) // 获取内容统计
	}
}



// healthCheck 健康检查接口
func healthCheck(c *gin.Context) {
	response.SuccessWithMessage(c, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
	}, "服务运行正常")
}
