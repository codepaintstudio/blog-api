package routes

import (
	"server/internal/controllers"
	"server/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// 使用中间件
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// 创建控制器实例
	userController := controllers.NewUserController(db)
	articleController := controllers.NewArticleController(db)

	// API 路由组
	api := router.Group("/api")

	// 用户相关路由
	user := api.Group("/user")
	{
		user.POST("/register", userController.Register)
		user.POST("/login", userController.Login)
	}

	// 文章相关路由
	article := api.Group("/articles")
	{
		// 公开路由
		article.GET("", articleController.List)         // 文章列表
		article.GET("/:id", articleController.GetById)  // 文章详情

		// 需要登录的路由
		auth := article.Group("", middleware.AuthMiddleware())
		{
			auth.POST("", articleController.Create)        // 创建文章
			auth.PUT("/:id", articleController.Update)      // 更新文章
			auth.DELETE("/:id", articleController.Delete)   // 删除文章
		}
	}
}
