package dragon

import (
	"valyria-backend/internal/apps/dragon/controller"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/service"
	"valyria-backend/internal/core/config"
)

type TektonTaskProvider struct {
	Repo       *repository.TektonTaskRepository
	Service    *service.TektonTaskManager
	Controller *controller.TektonTaskController
}

func NewTektonTaskProvider(cfg *config.Config) *TektonTaskProvider {
	repo := repository.NewTektonTaskRepository()
	svc := service.NewTektonTaskManager(repo, cfg)
	control := controller.NewTektonTaskController(svc)

	return &TektonTaskProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: control,
	}
}
