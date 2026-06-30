package prometheus

import (
	"prom-lens-backend/internal/apps/prometheus/controller"
	"prom-lens-backend/internal/apps/prometheus/repository"
	"prom-lens-backend/internal/apps/prometheus/service"

	"gorm.io/gorm"
)

type GroupProvider struct {
	Repo       *repository.GroupRepository
	Service    *service.GroupManager
	Controller *controller.GroupController
}

func NewGroupProvider(db *gorm.DB) *GroupProvider {
	repo := repository.NewGroupRepository(db)
	svc := service.NewGroupManager(repo)
	ctrl := controller.NewGroupController(svc)

	return &GroupProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
