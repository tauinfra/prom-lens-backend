package kingsguard

import (
	"gorm.io/gorm"
	"valyria-backend/internal/apps/kingsguard/controller"
	"valyria-backend/internal/apps/kingsguard/repository"
	"valyria-backend/internal/apps/kingsguard/service"
)

type PolicyProvider struct {
	Repo       *repository.PolicyRepository
	Service    *service.PolicyService
	Controller *controller.PolicyController
}

func NewPolicyProvider(db *gorm.DB) *PolicyProvider {
	// 初始化用户模块的依赖
	repo := repository.NewPolicyRepository(db)
	svc := service.NewPolicyService(repo, db)
	ctrl := controller.NewPolicyController(svc)

	return &PolicyProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
