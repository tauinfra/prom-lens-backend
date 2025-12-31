package auth

import (
	"valyria-backend/internal/apps/auth/controller"
	"valyria-backend/internal/apps/auth/repository"
	"valyria-backend/internal/apps/auth/service"

	"gorm.io/gorm"
)

type ModuleProvider struct {
	Repo       *repository.ModuleRepository
	Service    *service.ModuleManager
	Controller *controller.ModuleController
}

func NewModuleProvider(db *gorm.DB) *ModuleProvider {
	// 初始化用户模块的依赖
	repo := repository.NewModuleRepository(db)
	svc := service.NewModuleManager(db, repo)
	ctrl := controller.NewModuleController(svc)

	return &ModuleProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
