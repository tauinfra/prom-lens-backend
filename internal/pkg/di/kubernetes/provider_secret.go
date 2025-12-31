package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type SecretProvider struct {
	Repo       *repository.SecretRepository
	Service    *service.SecretManager
	Controller *controller.SecretController
}

func NewSecretProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *SecretProvider {
	// 初始化用户模块的依赖
	repo := repository.NewSecretRepository(cfgFactory, gvkFactory)
	svc := service.NewSecretManager(repo)
	ctrl := controller.NewSecretController(svc)

	return &SecretProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
