package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type DeploymentProvider struct {
	Repo       *repository.DeploymentRepository
	Service    *service.DeploymentService
	Controller *controller.DeploymentController
}

func NewDeploymentProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *DeploymentProvider {
	// 初始化用户模块的依赖
	repo := repository.NewDeploymentRepository(cfgFactory, gvkFactory)
	svc := service.NewDeploymentService(repo)
	ctrl := controller.NewDeploymentController(svc)

	return &DeploymentProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
