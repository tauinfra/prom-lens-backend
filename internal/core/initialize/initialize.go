package initialize

import (
	"fmt"
	"valyria-backend/internal/core/configs"
	"valyria-backend/internal/core/database"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/pkg/casbin/foundry"
	"valyria-backend/internal/pkg/casbin/k8s"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

var Encryptor encryption.Encryptor

func Components(cfg *configs.Config) (db *gorm.DB, encryptor encryption.Encryptor, err error) {
	// 初始化日志
	logger.InitLogger(cfg.Log.GetLoggerConfig())
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
	Encryptor = encryptor
	logger.Info("encryptor component initialized.")

	// 初始化 Redis
	if err = database.InitRedis(cfg); err != nil {
		return nil, nil, fmt.Errorf("redis init failed: %v", err)
	}

	// K8s Casbin 初始化
	if err = k8s.InitPolicy(db, cfg); err != nil {
		return nil, nil, err
	}
	// Foundry Casbin 初始化
	if err = foundry.InitPolicy(db, cfg); err != nil {
		return nil, nil, err
	}
	return db, encryptor, nil
}
