package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type PipelineProvider struct {
	Repo       *repository.PipelineRepository
	Service    *service.PipelineManager
	Controller *controller.PipelineController
}

func NewPipelineProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *PipelineProvider {
	// 初始化用户模块的依赖
	repo := repository.NewPipelineRepository(cfgFactory, gvkFactory)
	svc := service.NewPipelineManager(repo)
	ctrl := controller.NewPipelineController(svc)

	return &PipelineProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
