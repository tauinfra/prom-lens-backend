package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type DaemonSetProvider struct {
	Repo       *repository.DaemonSetRepository
	Service    *service.DaemonSetManager
	Controller *controller.DaemonSetController
}

func NewRDaemonSetProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *DaemonSetProvider {
	// 初始化用户模块的依赖
	repo := repository.NewDaemonSetRepository(cfgFactory, gvkFactory)
	svc := service.NewDaemonSetManager(repo)
	ctrl := controller.NewDaemonSetController(svc)

	return &DaemonSetProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
