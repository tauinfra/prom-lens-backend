package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type PersistentVolumeProvider struct {
	Repo       *repository.PersistentVolumeRepository
	Service    *service.PersistentVolumeManager
	Controller *controller.PersistentVolumeController
}

func NewPersistentVolumeProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *PersistentVolumeProvider {
	// 初始化用户模块的依赖
	repo := repository.NewPersistentVolumeRepository(cfgFactory, gvkFactory)
	svc := service.NewPersistentVolumeManager(repo)
	ctrl := controller.NewPersistentVolumeController(svc)

	return &PersistentVolumeProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
