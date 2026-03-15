package audit

import (
	"valyria-backend/internal/apps/audit/controller"
	"valyria-backend/internal/apps/audit/repository"
	"valyria-backend/internal/apps/audit/service"

	"gorm.io/gorm"
)

type AuditLogProvider struct {
	Repo       *repository.AuditLogRepository
	Service    *service.AuditLogService
	Controller *controller.AuditLogController
}

func NewAuditLogProvider(db *gorm.DB) *AuditLogProvider {
	repo := repository.NewAuditLogRepository(db)
	svc := service.NewAuditLogService(repo)
	ctrl := controller.NewAuditLogController(svc)

	return &AuditLogProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
