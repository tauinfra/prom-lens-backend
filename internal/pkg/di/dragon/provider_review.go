package dragon

import (
	"valyria-backend/internal/apps/dragon/controller"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/service"
	"valyria-backend/internal/core/config"

	"gorm.io/gorm"
)

type ReviewProvider struct {
	Repo       *repository.ReviewRepository
	Service    *service.ReviewManager
	Controller *controller.ReviewController
}

func NewReviewProvider(
	db *gorm.DB,
	configRepo repository.ReviewConfigRepository,
	aclRepo repository.PipelineACLRepository,
	gitlabManager service.GitlabManager,
	cfg *config.Config,
) *ReviewProvider {
	repo := repository.NewReviewRepository(db)
	svc := service.NewReviewManager(repo, configRepo, aclRepo, gitlabManager, db, cfg)
	ctrl := controller.NewReviewController(svc)

	return &ReviewProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
