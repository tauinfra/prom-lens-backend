package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/core/config"
	"valyria-backend/internal/pkg/k8s/factory"

	"gorm.io/gorm"
)

type PermissionProvider struct {
	Repo       repository.PermissionRepository
	Service    service.PermissionManager
	Controller *controller.PermissionController
}

// NewPermissionProvider 创建权限模块 Provider；sync 策略来自 config.permission_worker
func NewPermissionProvider(db *gorm.DB, cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory, cfg *config.Config) *PermissionProvider {
	permRepo := repository.NewPermissionRepository(db)
	clusterRoleBindingRepo := repository.NewClusterRoleBindingRepository(cfgFactory, gvkFactory)
	roleBindingRepo := repository.NewRoleBindingRepository(cfgFactory, gvkFactory)
	policy := service.SyncPolicy{}
	if cfg != nil {
		parsed := cfg.PermissionWorker.Parse()
		policy = service.SyncPolicy{
			MaxRetry:  parsed.MaxRetry,
			BaseDelay: parsed.BaseDelay,
			MaxDelay:  parsed.MaxDelay,
			BatchSize: parsed.BatchSize,
		}
	}
	svc := service.NewPermissionManager(permRepo, clusterRoleBindingRepo, roleBindingRepo, policy)
	ctrl := controller.NewPermissionController(svc)
	return &PermissionProvider{
		Repo:       permRepo,
		Service:    svc,
		Controller: ctrl,
	}
}
