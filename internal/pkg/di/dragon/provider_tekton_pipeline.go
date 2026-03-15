package dragon

import (
	"valyria-backend/internal/apps/dragon/controller"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/service"
	"valyria-backend/internal/core/config"
)

type TektonPipelineProvider struct {
	Repo       *repository.TektonPipelineRepository
	Service    *service.TektonPipelineManager
	Controller *controller.TektonPipelineController
}

func NewTektonPipelineProvider(cfg *config.Config) *TektonPipelineProvider {
	repo := repository.NewTektonPipelineRepository()
	svc := service.NewTektonPipelineManager(repo, cfg)
	control := controller.NewTektonPipelineController(svc)

	return &TektonPipelineProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: control,
	}
}
