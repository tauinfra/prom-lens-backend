package prometheus

import (
	"prom-lens-backend/internal/apps/prometheus/controller"
	"prom-lens-backend/internal/apps/prometheus/executor"
	"prom-lens-backend/internal/apps/prometheus/repository"
	"prom-lens-backend/internal/apps/prometheus/service"

	"gorm.io/gorm"
)

type RecordProvider struct {
	Repo       *repository.RecordRepository
	Service    *service.RecordManager
	Controller *controller.RecordController
}

func NewRecordProvider(db *gorm.DB, syncer executor.RuleSyncer) *RecordProvider {
	groupRepo := repository.NewGroupRepository(db)
	repo := repository.NewRecordRepository(db)
	svc := service.NewRecordManager(groupRepo, repo, syncer)
	ctrl := controller.NewRecordController(svc)

	return &RecordProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
