package initialize

import (
	"fmt"
	"valyria-backend/internal/core/config"
	"valyria-backend/internal/core/database"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

func Components(cfg *config.Config) (db *gorm.DB, encryptor encryption.Encryptor, err error) {
	// 初始化日志
	if err = logger.InitLogger(cfg.Log.GetLoggerConfig()); err != nil {
		return nil, nil, fmt.Errorf("logger init failed: %w", err)
	}
	logger.Info("Logger initialized")

	// 初始化数据库
	db, err = database.InitDatabase(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("database init failed: %w", err)
	}

	// 初始化加密器
	encryptor, err = encryption.NewSecretBoxEncryptor(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("create encryptor failed: %w", err)
	}
	logger.Info("encryptor component initialized.")

	// 初始化 Redis
	_, err = database.InitRedis(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("redis init failed: %w", err)
	}

	return db, encryptor, nil
}
