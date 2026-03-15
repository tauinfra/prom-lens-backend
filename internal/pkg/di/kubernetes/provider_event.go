package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
)

type EventProvider struct {
	Repo       *repository.EventRepository
	Service    *service.EventManager
	Controller *controller.EventController
}

func NewEventProvider(cfgFactory *repository.KubeConfigFactory) *EventProvider {
	// 初始化用户模块的依赖
	repo := repository.NewEventRepository(cfgFactory)
	svc := service.NewEventManager(repo)
	ctrl := controller.NewEventController(svc)

	return &EventProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
