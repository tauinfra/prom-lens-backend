package prometheus

import (
	"prom-lens-backend/internal/apps/prometheus/controller"
	"prom-lens-backend/internal/apps/prometheus/repository"
	"prom-lens-backend/internal/apps/prometheus/service"

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
