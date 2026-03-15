package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type RoleProvider struct {
	Repo       *repository.RoleRepository
	Service    *service.RoleManager
	Controller *controller.RoleController
}

func NewRoleProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *RoleProvider {
	repo := repository.NewRoleRepository(cfgFactory, gvkFactory)
	svc := service.NewRoleManager(repo)
	ctrl := controller.NewRoleController(svc)

	return &RoleProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
