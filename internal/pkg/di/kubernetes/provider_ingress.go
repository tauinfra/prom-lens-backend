package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type IngressProvider struct {
	Repo       *repository.IngressRepository
	Service    *service.IngressManager
	Controller *controller.IngressController
}

func NewIngressProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *IngressProvider {
	repo := repository.NewIngressRepository(cfgFactory, gvkFactory)
	svc := service.NewIngressManager(repo)
	ctrl := controller.NewIngressController(svc)

	return &IngressProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
