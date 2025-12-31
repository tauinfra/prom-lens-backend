package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type StatefulSetProvider struct {
	Repo       *repository.StatefulSetRepository
	Service    *service.StatefulSetService
	Controller *controller.StatefulSetController
}

func NewStatefulSetProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *StatefulSetProvider {
	// 初始化用户模块的依赖
	repo := repository.NewStatefulSetRepository(cfgFactory, gvkFactory)
	svc := service.NewStatefulSetService(repo)
	ctrl := controller.NewStatefulSetController(svc)

	return &StatefulSetProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
