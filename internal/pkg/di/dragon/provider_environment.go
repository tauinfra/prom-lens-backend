package dragon

import (
	"valyria-backend/internal/apps/dragon/controller"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/service"

	"gorm.io/gorm"
)

type EnvironmentProvider struct {
	Repo       *repository.EnvironmentRepository
	Service    *service.EnvironmentManager
	Controller *controller.EnvironmentController
}

func NewEnvironmentProvider(db *gorm.DB, aclRepo repository.PipelineACLRepository, credRepo repository.CredentialRepository) *EnvironmentProvider {
	repo := repository.NewEnvironmentRepository(db)
	svc := service.NewEnvironmentManager(repo, aclRepo, credRepo, db)
	control := controller.NewEnvironmentController(svc)

	return &EnvironmentProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: control,
	}
}
