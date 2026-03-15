package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type NamespaceProvider struct {
	Repo       *repository.NamespaceRepository
	Service    *service.NamespaceManager
	Controller *controller.NamespaceController
}

func NewNamespaceProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *NamespaceProvider {
	// 初始化用户模块的依赖
	repo := repository.NewNamespaceRepository(cfgFactory, gvkFactory)
	svc := service.NewNamespaceManager(repo)
	ctrl := controller.NewNamespaceController(svc)

	return &NamespaceProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
