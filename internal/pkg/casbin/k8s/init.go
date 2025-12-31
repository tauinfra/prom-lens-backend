package k8s

import (
	"valyria-backend/internal/core/configs"
	"valyria-backend/internal/core/logger"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var Enforcer *casbin.Enforcer

func InitPolicy(tx *gorm.DB, cfg *configs.Config) error {
	// 使用 Gorm 适配器(自定义表名)
	adapter, err := gormadapter.NewAdapterByDBUseTableName(tx, "", "valyria_kubernetes_policy")
	if err != nil {
		logger.Fatalf("Kubernetes module Casbin Adapter initialization failed. err: %v", err.Error())
		return err
	}
	// 创建执行器，加载模型文件
	e, err := casbin.NewEnforcer("internal/pkg/casbin/k8s/model.conf", adapter)
	if err != nil {
		logger.Fatalf("Kubernetes module Casbin Enforcer initialization failed. err: %v", err.Error())
		return err
	}

	// 从数据库加载策略
	if err = e.LoadPolicy(); err != nil {
		logger.Fatalf("Kubernetes module Casbin failed to load permission k8s. err: %v", err.Error())
		return err
	}

	// 启用自动保存，Web 操作自动持久化
	e.EnableAutoSave(true)

	Enforcer = e

	logger.Info("Kubernetes module Casbin initialized successfully.")
	return nil
}
