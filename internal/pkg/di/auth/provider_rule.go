package auth

import (
	"valyria-backend/internal/apps/auth/controller"
	"valyria-backend/internal/apps/auth/repository"
	"valyria-backend/internal/apps/auth/service"

	"gorm.io/gorm"
)

type RuleProvider struct {
	Repo       *repository.RoleRepository
	Service    *service.RoleManager
	Controller *controller.RoleController
}

func NewRuleProvider(db *gorm.DB) *RuleProvider {
	// 初始化用户模块的依赖
	repo := repository.NewRoleRepository(db)
	svc := service.NewRoleManager(db, repo)
	ctrl := controller.NewRoleController(svc)

	return &RuleProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
