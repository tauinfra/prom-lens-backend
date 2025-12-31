package kingsguard

import (
	"valyria-backend/internal/core/configs"
	"valyria-backend/internal/core/logger"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var Enforcer *casbin.Enforcer

func InitCasbin(tx *gorm.DB, cfg *configs.Config) error {
	// 使用 Gorm 适配器(自定义表名)
	adapter, err := gormadapter.NewAdapterByDBUseTableName(tx, "", cfg.Policy.PolicyTable)
	if err != nil {
		logger.Fatalf("Casbin Adapter 初始化失败: %v", err.Error())
		return err
	}
	// 创建执行器，加载模型文件
	e, err := casbin.NewEnforcer(cfg.Policy.ModelPath, adapter)
	if err != nil {
		logger.Fatalf("Casbin Enforcer 初始化失败: %v", err.Error())
		return err
	}

	// 从数据库加载策略
	if err = e.LoadPolicy(); err != nil {
		logger.Fatalf("Casbin 加载策略失败: %v", err.Error())
		return err
	}

	// 启用自动保存，Web 操作自动持久化
	e.EnableAutoSave(cfg.Policy.AutoSave)

	Enforcer = e

	logger.Info("Casbin 权限系统初始化成功")
	return nil
}
