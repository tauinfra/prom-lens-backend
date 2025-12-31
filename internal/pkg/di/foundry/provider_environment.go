package foundry

import (
	"valyria-backend/internal/apps/foundry/controller"
	"valyria-backend/internal/apps/foundry/repository"
	"valyria-backend/internal/apps/foundry/service"

	"gorm.io/gorm"
)

type EnvironmentProvider struct {
	Repo       *repository.EnvironmentRepository
	Service    *service.EnvironmentManager
	Controller *controller.EnvironmentController
}

func NewEnvironmentProvider(db *gorm.DB) *EnvironmentProvider {
	repo := repository.NewEnvironmentRepository(db)
	svc := service.NewEnvironmentManager(repo, db)
	ctrl := controller.NewEnvironmentController(svc)

	return &EnvironmentProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
