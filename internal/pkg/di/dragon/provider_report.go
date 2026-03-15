package dragon

import (
	"valyria-backend/internal/apps/dragon/controller"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/service"

	"gorm.io/gorm"
)

type ReportProvider struct {
	Controller *controller.ReleaseReportController
}

func NewReportProvider(db *gorm.DB) *ReportProvider {
	reportRepo := repository.NewReleaseReportRepository(db)
	reportSvc := service.NewReleaseReportManager(reportRepo)
	reportCtrl := controller.NewReleaseReportController(reportSvc)
	return &ReportProvider{
		Controller: reportCtrl,
	}
}
