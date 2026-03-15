package audit

import (
	"valyria-backend/internal/apps/audit/controller"
	"valyria-backend/internal/apps/audit/repository"
	"valyria-backend/internal/apps/audit/service"

	"gorm.io/gorm"
)

type AuthLogProvider struct {
	Repo       *repository.AuthLogRepository
	Service    *service.AuthLogService
	Controller *controller.AuthLogController
}

func NewAuthLogProvider(db *gorm.DB) *AuthLogProvider {
	repo := repository.NewAuthLogRepository(db)
	svc := service.NewAuthLogService(repo)
	ctrl := controller.NewAuthLogController(svc)

	return &AuthLogProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
