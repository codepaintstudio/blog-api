package main

import (
	"fmt"
	"os"
	"server/internal/models"
	"server/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// 从环境变量中获取数据库连接信息
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移表结构
	db.AutoMigrate(&models.User{}, &models.Article{})

	// 初始化路由
	router := gin.Default()

	// 设置代理
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// 设置路由
	routes.SetupRoutes(router, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "7777"
	}

	// 启动服务
	if err := router.Run(":" + port); err != nil {
		panic("failed to start server")
	} else {
		fmt.Println("Server started on :" + port)
	}
}
