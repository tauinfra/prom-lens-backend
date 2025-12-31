package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type ServiceProvider struct {
	Repo       *repository.ServiceRepository
	Service    *service.ServiceManager
	Controller *controller.ServiceController
}

func NewServiceProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *ServiceProvider {
	// 初始化用户模块的依赖
	repo := repository.NewServiceRepository(cfgFactory, gvkFactory)
	svc := service.NewServiceManager(repo)
	ctrl := controller.NewServiceController(svc)

	return &ServiceProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
