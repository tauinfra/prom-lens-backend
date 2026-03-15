package prometheus

import (
	"valyria-backend/internal/apps/prometheus/controller"
	"valyria-backend/internal/apps/prometheus/executor"
	"valyria-backend/internal/apps/prometheus/repository"
	"valyria-backend/internal/apps/prometheus/service"

	"gorm.io/gorm"
)

type RecordProvider struct {
	Repo       *repository.RecordRepository
	Service    *service.RecordManager
	Controller *controller.RecordController
}

func NewRecordProvider(db *gorm.DB, syncer executor.RuleSyncer) *RecordProvider {
	repo := repository.NewRecordRepository(db)
	svc := service.NewRecordManager(repo, syncer)
	ctrl := controller.NewRecordController(svc)

	return &RecordProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
