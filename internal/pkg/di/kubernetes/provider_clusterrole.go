package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type ClusterRoleProvider struct {
	Repo       *repository.ClusterRoleRepository
	Service    *service.ClusterRoleManager
	Controller *controller.ClusterRoleController
}

func NewClusterRoleProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *ClusterRoleProvider {
	repo := repository.NewClusterRoleRepository(cfgFactory, gvkFactory)
	svc := service.NewClusterRoleManager(repo)
	ctrl := controller.NewClusterRoleController(svc)

	return &ClusterRoleProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
