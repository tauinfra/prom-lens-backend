package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type PersistentVolumeClaimProvider struct {
	Repo       *repository.PersistentVolumeClaimRepository
	Service    *service.PersistentVolumeClaimManager
	Controller *controller.PersistentVolumeClaimController
}

func NewPersistentVolumeClaimProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *PersistentVolumeClaimProvider {
	// 初始化用户模块的依赖
	repo := repository.NewPersistentVolumeClaimRepository(cfgFactory, gvkFactory)
	svc := service.NewPersistentVolumeClaimManager(repo)
	ctrl := controller.NewPersistentVolumeClaimController(svc)

	return &PersistentVolumeClaimProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
