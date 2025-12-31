package foundry

import (
	"valyria-backend/internal/apps/foundry/controller"
	"valyria-backend/internal/apps/foundry/repository"
	"valyria-backend/internal/apps/foundry/service"

	"gorm.io/gorm"
)

type ApplicationProvider struct {
	Repo       *repository.ApplicationRepository
	Service    *service.ApplicationManager
	Controller *controller.ApplicationController
}

func NewApplicationProvider(db *gorm.DB) *ApplicationProvider {
	repo := repository.NewApplicationRepository(db)
	svc := service.NewApplicationManager(repo, db)
	ctrl := controller.NewApplicationController(svc)

	return &ApplicationProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}

type AppContext struct {
	Gitlab  *GitlabProvider
	Release *ReleaseProvider
}
