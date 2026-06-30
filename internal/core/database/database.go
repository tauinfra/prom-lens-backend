package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"prom-lens-backend/internal/core/config"
	"prom-lens-backend/internal/core/logger"
)

// InitDatabase 非全局
func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	// 获取通用数据库对象 sql.DB 以使用其提供的功能
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)       // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)       // 设置数据库的最大打开连接数
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime) // 设置连接的最大可复用时间
	logger.Info("MySQL database connected successfully.")
	return db, nil
}

// CloseDatabase 关闭数据库连接
func CloseDatabase(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	logger.Info("MySQL database connection closed.")
	return nil
}
