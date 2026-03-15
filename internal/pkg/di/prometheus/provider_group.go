package prometheus

import (
	"valyria-backend/internal/apps/prometheus/controller"
	"valyria-backend/internal/apps/prometheus/repository"
	"valyria-backend/internal/apps/prometheus/service"

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
