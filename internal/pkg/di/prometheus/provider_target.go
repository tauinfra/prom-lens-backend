package prometheus

import (
	"prom-lens-backend/internal/apps/prometheus/controller"
	"prom-lens-backend/internal/apps/prometheus/executor"
	"prom-lens-backend/internal/apps/prometheus/repository"
	"prom-lens-backend/internal/apps/prometheus/service"

	"gorm.io/gorm"
)

type TargetProvider struct {
	Repo       *repository.TargetRepository
	Service    *service.TargetManager
	Controller *controller.TargetController
}

func NewTargetProvider(db *gorm.DB, syncer executor.TargetSyncer) *TargetProvider {
	repo := repository.NewTargetRepository(db)
	svc := service.NewTargetManager(repo, syncer)
	ctrl := controller.NewTargetController(svc)

	return &TargetProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
