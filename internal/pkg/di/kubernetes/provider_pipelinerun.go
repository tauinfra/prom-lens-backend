package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type PipelineRunProvider struct {
	Repo       *repository.PipelineRunRepository
	Service    *service.PipelineRunManager
	Controller *controller.PipelineRunController
}

func NewPipelineRunProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *PipelineRunProvider {
	// 初始化用户模块的依赖
	repo := repository.NewPipelineRunRepository(cfgFactory, gvkFactory)
	svc := service.NewPipelineRunManager(repo)
	ctrl := controller.NewPipelineRunController(svc)

	return &PipelineRunProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
