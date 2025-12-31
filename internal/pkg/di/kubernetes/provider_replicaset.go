package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type ReplicaSetProvider struct {
	Repo       *repository.ReplicaSetRepository
	Service    *service.ReplicaSetManager
	Controller *controller.ReplicaSetController
}

func NewReplicaSetProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *ReplicaSetProvider {
	// 初始化用户模块的依赖
	repo := repository.NewReplicaSetRepository(cfgFactory, gvkFactory)
	svc := service.NewReplicaSetManager(repo)
	ctrl := controller.NewReplicaSetController(svc)

	return &ReplicaSetProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
