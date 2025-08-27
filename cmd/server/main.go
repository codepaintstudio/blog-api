// Package main Blog API
// @title Blog API
// @version 1.0
// @description 一个专为前端开发学习者提供的博客后端API系统
// @termsOfService https://github.com/codepaintstudio/blog-api
//
// @contact.name API Support
// @contact.url https://github.com/codepaintstudio/blog-api/issues
// @contact.email support@example.com
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"blog-api/internal/container"
	"blog-api/internal/routes"
	"blog-api/pkg/config"
	"blog-api/pkg/database"
	"blog-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.Info("博客服务启动中...")

	// 初始化配置
	cfg, err := config.LoadConfig("")
	if err != nil {
		logger.Fatal("加载配置失败", logger.Err("error", err))
	}
	logger.Info("配置加载成功")

	// 初始化日志系统
	if err := logger.InitLogger(&cfg.Log); err != nil {
		logger.Fatal("初始化日志系统失败", logger.Err("error", err))
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库连接
	if err := database.InitMySQL(); err != nil {
		logger.Fatal("MySQL 初始化失败", logger.Err("error", err))
	}

	if err := database.InitRedis(); err != nil {
		logger.Fatal("Redis 初始化失败", logger.Err("error", err))
	}

	// 初始化依赖注入容器
	appContainer := container.NewContainer()
	appContainer.InitRepositories()
	appContainer.InitServices()
	appContainer.InitControllers()

	// 初始化系统管理员
	adminInitService := appContainer.GetAdminInitService()
	if err := adminInitService.InitializeAdmin(); err != nil {
		logger.Error("管理员初始化失败", logger.Err("error", err))
	}

	// 获取所有控制器
	authController, userController, articleController, categoryController, fileController, commentController, likeController, favoriteController, adminController := appContainer.GetControllers()

	// 初始化路由
	router := routes.SetupRoutes(authController, userController, articleController, categoryController, fileController, commentController, likeController, favoriteController, adminController)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// 设置优雅关闭
	go setupGracefulShutdown(srv, appContainer)

	logger.Info("博客API服务器已启动",
		logger.String("address", fmt.Sprintf("http://localhost:%d", cfg.Server.Port)),
		logger.String("mode", cfg.Server.Mode),
		logger.String("swagger", fmt.Sprintf("http://localhost:%d/swagger/index.html", cfg.Server.Port)),
	)

	// 启动服务器
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("服务器启动失败", logger.Err("error", err))
	}
}

func setupGracefulShutdown(srv *http.Server, appContainer *container.Container) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Info("正在关闭服务器...")

		// 设置关闭超时
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// 关闭HTTP服务器
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("关闭HTTP服务器时出错", logger.Err("error", err))
		}

		// 关闭数据库连接
		if err := database.CloseMySQL(); err != nil {
			logger.Error("关闭MySQL时出错", logger.Err("error", err))
		}

		if err := database.CloseRedis(); err != nil {
			logger.Error("关闭Redis时出错", logger.Err("error", err))
		}

		logger.Info("服务器关闭完成")
		os.Exit(0)
	}()
}
