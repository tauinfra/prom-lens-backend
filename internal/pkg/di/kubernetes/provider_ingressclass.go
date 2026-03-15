package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type IngressClassProvider struct {
	Repo       *repository.IngressClassRepository
	Service    *service.IngressClassManager
	Controller *controller.IngressClassController
}

func NewIngressClassProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *IngressClassProvider {
	repo := repository.NewIngressClassRepository(cfgFactory, gvkFactory)
	svc := service.NewIngressClassManager(repo)
	ctrl := controller.NewIngressClassController(svc)

	return &IngressClassProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
