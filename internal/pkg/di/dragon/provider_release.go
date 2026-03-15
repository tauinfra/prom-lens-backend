package dragon

import (
	"valyria-backend/internal/apps/dragon/controller"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/service"
	"valyria-backend/internal/core/config"

	"gorm.io/gorm"
)

type ReleaseProvider struct {
	Repo       *repository.ReleaseRepository
	Service    *service.ReleaseManager
	Controller *controller.ReleaseController
}

func NewReleaseProvider(
	db *gorm.DB,
	pipelineRepo repository.PipelineRepository,
	aclRepo repository.PipelineACLRepository,
	reviewManager service.ReviewManager,
	gitlabManager service.GitlabManager,
	cfg *config.Config,
) *ReleaseProvider {
	repo := repository.NewReleaseRepository(db)
	svc := service.NewReleaseManager(repo, pipelineRepo, aclRepo, reviewManager, gitlabManager, db, cfg)
	control := controller.NewReleaseController(svc)

	return &ReleaseProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: control,
	}
}
