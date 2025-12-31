package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type TaskRunProvider struct {
	Repo       *repository.TaskRunRepository
	Service    *service.TaskRunManager
	Controller *controller.TaskRunController
}

func NewTaskRunProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *TaskRunProvider {
	// 初始化用户模块的依赖
	repo := repository.NewTaskRunRepository(cfgFactory, gvkFactory)
	svc := service.NewTaskRunManager(repo)
	ctrl := controller.NewTaskRunController(svc)

	return &TaskRunProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
