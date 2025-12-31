package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/k8s/factory"
)

type TaskProvider struct {
	Repo       *repository.TaskRepository
	Service    *service.TaskManager
	Controller *controller.TaskController
}

func NewTaskProvider(cfgFactory *repository.KubeConfigFactory, gvkFactory *factory.GVKFactory) *TaskProvider {
	// 初始化用户模块的依赖
	repo := repository.NewTaskRepository(cfgFactory, gvkFactory)
	svc := service.NewTaskManager(repo)
	ctrl := controller.NewTaskController(svc)

	return &TaskProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
