package dragon

import (
	"valyria-backend/internal/apps/dragon/controller"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/service"

	"gorm.io/gorm"
)

type PipelineProvider struct {
	Repo       *repository.PipelineRepository
	Service    *service.PipelineManager
	Controller *controller.PipelineController
}

func NewPipelineProvider(db *gorm.DB, credentialRepo repository.CredentialRepository, aclRepo repository.PipelineACLRepository) *PipelineProvider {
	repo := repository.NewPipelineRepository(db)
	svc := service.NewPipelineManager(repo, aclRepo, db)
	control := controller.NewPipelineController(svc)

	return &PipelineProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: control,
	}
}
