package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"valyria-backend/internal/core/configs"
	"valyria-backend/internal/core/logger"
)

// InitDatabase 非全局
func InitDatabase(cfg *configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("MySQL database connected failure, %v\n", err)
		return nil, err
	}
	// 获取通用数据库对象 sql.DB 以使用其提供的功能
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatalf("MySQL database et sql.DB failed, %v\\", err)
		return nil, err
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)       // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)       // 设置数据库的最大打开连接数
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime) // 设置连接的最大可复用时间
	logger.Info("MySQL database connected successfully.")
	return db, nil
}
