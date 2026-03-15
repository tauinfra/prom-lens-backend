package dragon

import (
	"valyria-backend/internal/apps/dragon/controller"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/service"

	"gorm.io/gorm"
)

type ProjectProvider struct {
	Repo       *repository.ProjectRepository
	Service    *service.ProjectManager
	Controller *controller.ProjectController
}

func NewProjectProvider(db *gorm.DB, aclRepo repository.PipelineACLRepository) *ProjectProvider {
	repo := repository.NewProjectRepository(db)
	svc := service.NewProjectManager(repo, aclRepo, db)
	control := controller.NewProjectController(svc)

	return &ProjectProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: control,
	}
}
