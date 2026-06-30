package initialize

import (
	"fmt"
	"prom-lens-backend/internal/core/config"
	"prom-lens-backend/internal/core/database"
	"prom-lens-backend/internal/core/logger"

	"gorm.io/gorm"
)

func Components(cfg *config.Config) (db *gorm.DB, err error) {
	if err = logger.InitLogger(cfg.Log.GetLoggerConfig()); err != nil {
		return nil, fmt.Errorf("logger init failed: %w", err)
	}
	logger.Info("Logger initialized")

	db, err = database.InitDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("database init failed: %w", err)
	}

	return db, nil
}
