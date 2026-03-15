package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type ConfigmapProvider struct {
	Repo       *repository.ConfigmapRepository
	Service    *service.ConfigmapManager
	Controller *controller.ConfigmapController
}

func NewConfigmapProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *ConfigmapProvider {
	// 初始化用户模块的依赖
	repo := repository.NewConfigmapRepository(cfgFactory, gvkFactory)
	svc := service.NewConfigmapManager(repo)
	ctrl := controller.NewConfigmapController(svc)

	return &ConfigmapProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
