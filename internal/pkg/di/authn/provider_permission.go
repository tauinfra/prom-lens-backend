package authn

import (
	"valyria-backend/internal/apps/authn/controller"
	"valyria-backend/internal/apps/authn/repository"
	"valyria-backend/internal/apps/authn/service"

	"gorm.io/gorm"
)

type PermissionProvider struct {
	Repo       *repository.PermissionRepository
	Service    *service.PermissionManager
	Controller *controller.PermissionController
}

func NewPermissionProvider(db *gorm.DB) *PermissionProvider {
	// 初始化用户模块的依赖
	repo := repository.NewPermissionRepository(db)
	manager := service.NewPermissionManager(db, repo)
	control := controller.NewPermissionController(manager)

	return &PermissionProvider{
		Repo:       &repo,
		Service:    &manager,
		Controller: control,
	}
}
