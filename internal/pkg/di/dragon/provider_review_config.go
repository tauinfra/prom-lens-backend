package dragon

import (
	"valyria-backend/internal/apps/dragon/controller"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/service"

	"gorm.io/gorm"
)

type ReviewConfigProvider struct {
	Repo       *repository.ReviewConfigRepository
	Service    *service.ReviewConfigManager
	Controller *controller.ReviewConfigController
}

func NewReviewConfigProvider(
	db *gorm.DB,
) *ReviewConfigProvider {
	repo := repository.NewReviewConfigRepository(db)
	svc := service.NewReviewConfigManager(repo, db)
	ctrl := controller.NewReviewConfigController(svc)

	return &ReviewConfigProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
