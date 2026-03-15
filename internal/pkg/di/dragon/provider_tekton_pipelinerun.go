package dragon

import (
	"valyria-backend/internal/apps/dragon/controller"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/service"
	"valyria-backend/internal/core/config"
)

type TektonPipelineRunProvider struct {
	Repo       *repository.TektonPipelineRunRepository
	Service    *service.TektonPipelineRunManager
	Controller *controller.TektonPipelineRunController
}

func NewTektonPipelineRunProvider(cfg *config.Config) *TektonPipelineRunProvider {
	repo := repository.NewTektonPipelineRunRepository()
	svc := service.NewTektonPipelineRunManager(repo, cfg)
	control := controller.NewTektonPipelineRunController(svc)

	return &TektonPipelineRunProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: control,
	}
}
