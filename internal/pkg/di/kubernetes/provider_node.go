package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type NodeProvider struct {
	Repo       *repository.NodeRepository
	Service    *service.NodeService
	Controller *controller.NodeController
}

func NewNodeProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *NodeProvider {
	// 初始化用户模块的依赖
	repo := repository.NewNodeRepository(cfgFactory, gvkFactory)
	svc := service.NewNodeService(repo)
	ctrl := controller.NewNodeController(svc)

	return &NodeProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
