package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"valyria-backend/internal/core/config"
	"valyria-backend/internal/core/logger"
)

var Redis *redis.Client // Deprecated: 使用 InitRedis 返回的客户端，全局变量仅用于向后兼容

func InitRedis(cfg *config.Config) (*redis.Client, error) {
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
		return nil, fmt.Errorf("failed to connect redis %s: %w", cfg.Redis.Addr, err)
	}

	// 设置全局变量（向后兼容）
	Redis = rdb

	logger.Infof("Redis connected successfully: %s", cfg.Redis.Addr)
	return rdb, nil
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() error {
	if Redis == nil {
		return nil
	}
	if err := Redis.Close(); err != nil {
		return fmt.Errorf("failed to close redis: %w", err)
	}
	logger.Info("Redis connection closed.")
	return nil
}
