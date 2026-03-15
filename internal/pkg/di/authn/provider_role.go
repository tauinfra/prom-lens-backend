package authn

import (
	"valyria-backend/internal/apps/authn/controller"
	"valyria-backend/internal/apps/authn/repository"
	"valyria-backend/internal/apps/authn/service"

	"gorm.io/gorm"
)

type RoleProvider struct {
	Repo       *repository.RoleRepository
	Service    *service.RoleManager
	Controller *controller.RoleController
}

func NewRoleProvider(db *gorm.DB) *RoleProvider {
	// 初始化用户模块的依赖
	repo := repository.NewRoleRepository(db)
	manager := service.NewRoleManager(db, repo)
	control := controller.NewRoleController(manager)

	return &RoleProvider{
		Repo:       &repo,
		Service:    &manager,
		Controller: control,
	}
}
