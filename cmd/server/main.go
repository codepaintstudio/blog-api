package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"blog-api/pkg/config"
	"blog-api/pkg/database"
)

func main() {
	log.Println("博客服务启动中...")

	// 初始化配置
	_, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("加载环境变量失败: %v", err)
	}
	log.Println("配置加载成功")

	// 初始化数据库连接
	if err := database.InitMySQL(); err != nil {
		log.Fatalf("MYSQL 初始化失败: %v", err)
	}

	if err := database.InitRedis(); err != nil {
		log.Fatalf("REDIS 初始化失败: %v", err)
	}

	// 设置优雅关闭
	setupGracefulShutdown()

	// TODO: 初始化路由和中间件
	// TODO: 启动服务器

	log.Println("博客API服务器已就绪")

	// 等待关闭信号
	waitForShutdown()
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("正在关闭服务器...")

		// 关闭数据库连接
		if err := database.CloseMySQL(); err != nil {
			log.Printf("关闭MySQL时出错: %v", err)
		}

		if err := database.CloseRedis(); err != nil {
			log.Printf("关闭Redis时出错: %v", err)
		}

		log.Println("服务器关闭完成")
		os.Exit(0)
	}()
}

func waitForShutdown() {
	select {}
}