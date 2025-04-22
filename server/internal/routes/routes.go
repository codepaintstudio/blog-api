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

	// 帖子相关路由
	article := api.Group("/articles")
	{
		// 公开路由
		article.GET("", articleController.List)        // 帖子列表
		article.GET("/:id", articleController.GetById) // 帖子详情

		// 需要登录的路由
		auth := article.Group("", middleware.AuthMiddleware())
		{
			auth.POST("", articleController.Create)       // 创建帖子
			auth.PUT("/:id", articleController.Update)    // 更新帖子
			auth.DELETE("/:id", articleController.Delete) // 删除帖子
		}
	}
}
