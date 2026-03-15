package dragon

import (
	"valyria-backend/internal/apps/dragon/controller"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/service"

	"gorm.io/gorm"
)

type PipelineACLProvider struct {
	Repo       *repository.PipelineACLRepository
	Service    *service.PipelineACLManager
	Controller *controller.PipelineACLController
}

func NewPipelineACLProvider(db *gorm.DB) *PipelineACLProvider {
	repo := repository.NewPipelineACLRepository(db)
	svc := service.NewPipelineACLManager(repo)
	ctrl := controller.NewPipelineACLController(svc)

	return &PipelineACLProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
