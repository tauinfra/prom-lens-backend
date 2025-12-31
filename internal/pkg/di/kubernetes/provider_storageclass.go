package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type StorageClassProvider struct {
	Repo       *repository.StorageClassRepository
	Service    *service.StorageClassManager
	Controller *controller.StorageClassController
}

func NewStorageClassProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *StorageClassProvider {
	// 初始化用户模块的依赖
	repo := repository.NewStorageClassRepository(cfgFactory, gvkFactory)
	svc := service.NewStorageClassManager(repo)
	ctrl := controller.NewStorageClassController(svc)

	return &StorageClassProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
