package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type HpaProvider struct {
	Repo       *repository.HpaRepository
	Service    *service.HPAManager
	Controller *controller.HPAController
}

func NewHpaProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory, hpaHistory service.HpaHistoryManager) *HpaProvider {
	repo := repository.NewHpaRepository(cfgFactory, gvkFactory)
	svc := service.NewHPAManager(repo)
	ctrl := controller.NewHPAController(svc, hpaHistory)
	return &HpaProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
