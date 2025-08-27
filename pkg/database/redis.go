package database

import (
	"context"
	"fmt"
	"log"

	"blog-api/pkg/config"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() error {
	cfg := config.GetConfig()
	if cfg == nil {
		return fmt.Errorf("配置未加载")
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("连接Redis失败: %w", err)
	}

	log.Println("Redis连接成功")
	return nil
}

func GetRedis() *redis.Client {
	return RedisClient
}

func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}