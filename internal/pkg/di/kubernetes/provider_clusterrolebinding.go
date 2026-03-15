package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type ClusterRoleBindingProvider struct {
	Repo       *repository.ClusterRoleBindingRepository
	Service    *service.ClusterRoleBindingManager
	Controller *controller.ClusterRoleBindingController
}

func NewClusterRoleBindingProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *ClusterRoleBindingProvider {
	repo := repository.NewClusterRoleBindingRepository(cfgFactory, gvkFactory)
	svc := service.NewClusterRoleBindingManager(repo)
	ctrl := controller.NewClusterRoleBindingController(svc)

	return &ClusterRoleBindingProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
