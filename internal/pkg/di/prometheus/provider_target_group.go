package prometheus

import (
	"valyria-backend/internal/apps/prometheus/controller"
	"valyria-backend/internal/apps/prometheus/repository"
	"valyria-backend/internal/apps/prometheus/service"

	"gorm.io/gorm"
)

type TargetGroupProvider struct {
	Repo       *repository.TargetGroupRepository
	Service    *service.TargetGroupManager
	Controller *controller.TargetGroupController
}

func NewTargetGroupProvider(db *gorm.DB) *TargetGroupProvider {
	repo := repository.NewTargetGroupRepository(db)
	svc := service.NewTargetGroupManager(repo)
	ctrl := controller.NewTargetGroupController(svc)

	return &TargetGroupProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
