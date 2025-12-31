package foundry

import (
	"valyria-backend/internal/apps/foundry/controller"
	"valyria-backend/internal/apps/foundry/repository"
	"valyria-backend/internal/apps/foundry/service"

	"gorm.io/gorm"
)

type ReleaseProvider struct {
	Repo       *repository.ReleaseRepository
	Service    *service.ReleaseManager
	Controller *controller.ReleaseController
}

func NewReleaseProvider(db *gorm.DB, pipelineRepository repository.PipelineRepository, gitlabManager service.GitlabManager) *ReleaseProvider {

	repo := repository.NewReleaseRepository(db)
	svc := service.NewReleaseManager(repo, pipelineRepository, gitlabManager, db)
	ctrl := controller.NewReleaseController(svc)

	return &ReleaseProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
