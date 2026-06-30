package prometheus

import (
	"prom-lens-backend/internal/apps/prometheus/controller"
	"prom-lens-backend/internal/apps/prometheus/repository"
	"prom-lens-backend/internal/apps/prometheus/service"

	"gorm.io/gorm"
)

type SyncProvider struct {
	Service    service.SyncManager
	Controller *controller.SyncController
}

func NewSyncProvider(db *gorm.DB) *SyncProvider {
	groupRepo := repository.NewGroupRepository(db)
	ruleRepo := repository.NewRuleRepository(db)
	recordRepo := repository.NewRecordRepository(db)
	svc := service.NewSyncManager(groupRepo, ruleRepo, recordRepo)
	ctrl := controller.NewSyncController(svc)
	return &SyncProvider{
		Service:    svc,
		Controller: ctrl,
	}
}
