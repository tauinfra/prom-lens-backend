package foundry

import (
	"valyria-backend/internal/apps/foundry/controller"
	"valyria-backend/internal/apps/foundry/repository"
	"valyria-backend/internal/apps/foundry/service"

	"gorm.io/gorm"
)

type PipelineProvider struct {
	Repo       *repository.PipelineRepository
	Service    *service.PipelineManager
	Controller *controller.PipelineController
}

func NewPipelineProvider(db *gorm.DB) *PipelineProvider {
	repo := repository.NewPipelineRepository(db)
	svc := service.NewPipelineManager(repo, db)
	ctrl := controller.NewPipelineController(svc)

	return &PipelineProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
