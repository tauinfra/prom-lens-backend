package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"valyria-backend/internal/core/configs"
	"valyria-backend/internal/core/logger"
)

var Redis *redis.Client

func InitRedis(cfg *configs.Config) error {
	// 设置默认值
	poolSize := cfg.Redis.PoolSize
	if poolSize == 0 {
		poolSize = 10 // 默认连接池大小
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     poolSize,
		MinIdleConns: 5, // 最小空闲连接数
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试连接
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect redis %s: %w", cfg.Redis.Addr, err)
	}

	// 设置全局变量
	Redis = rdb

	logger.Infof("Redis connected successfully: %s", cfg.Redis.Addr)
	return nil
}
