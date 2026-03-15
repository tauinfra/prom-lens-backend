package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type PodProvider struct {
	Repo       *repository.PodRepository
	Service    *service.PodManager
	Controller *controller.PodController
}

func NewPodProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *PodProvider {
	// 初始化用户模块的依赖
	repo := repository.NewPodRepository(cfgFactory, gvkFactory)
	svc := service.NewPodManager(repo)
	ctrl := controller.NewPodController(svc)

	return &PodProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
