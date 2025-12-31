package foundry

import (
	"valyria-backend/internal/apps/foundry/controller"
	"valyria-backend/internal/apps/foundry/repository"
	"valyria-backend/internal/apps/foundry/service"

	"gorm.io/gorm"
)

type ProjectProvider struct {
	Repo       *repository.ProjectRepository
	Service    *service.ProjectManager
	Controller *controller.ProjectController
}

func NewProjectProvider(db *gorm.DB) *ProjectProvider {
	repo := repository.NewProjectRepository(db)
	svc := service.NewProjectManager(repo, db)
	ctrl := controller.NewProjectController(svc)

	return &ProjectProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
