package dragon

import (
	"valyria-backend/internal/apps/dragon/controller"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/service"
	"valyria-backend/internal/core/config"
)

type TektonTaskRunProvider struct {
	Repo       *repository.TektonTaskRunRepository
	Service    *service.TektonTaskRunManager
	Controller *controller.TektonTaskRunController
}

func NewTektonTaskRunProvider(cfg *config.Config) *TektonTaskRunProvider {
	repo := repository.NewTektonTaskRunRepository()
	svc := service.NewTektonTaskRunManager(repo, cfg)
	control := controller.NewTektonTaskRunController(svc)

	return &TektonTaskRunProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: control,
	}
}
