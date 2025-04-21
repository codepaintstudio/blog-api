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

	// API 路由组
	api := router.Group("/api")

	// 用户相关路由
	user := api.Group("/user")
	{
		user.POST("/register", userController.Register)
		user.POST("/login", userController.Login)
	}
}
