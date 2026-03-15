package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type ServiceAccountProvider struct {
	Repo       *repository.ServiceAccountRepository
	Service    *service.ServiceAccountManager
	Controller *controller.ServiceAccountController
}

func NewServiceAccountProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *ServiceAccountProvider {
	repo := repository.NewServiceAccountRepository(cfgFactory, gvkFactory)
	svc := service.NewServiceAccountManager(repo)
	ctrl := controller.NewServiceAccountController(svc)

	return &ServiceAccountProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
