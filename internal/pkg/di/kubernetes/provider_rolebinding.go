package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type RoleBindingProvider struct {
	Repo       *repository.RoleBindingRepository
	Service    *service.RoleBindingManager
	Controller *controller.RoleBindingController
}

func NewRoleBindingProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *RoleBindingProvider {
	repo := repository.NewRoleBindingRepository(cfgFactory, gvkFactory)
	svc := service.NewRoleBindingManager(repo)
	ctrl := controller.NewRoleBindingController(svc)

	return &RoleBindingProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
