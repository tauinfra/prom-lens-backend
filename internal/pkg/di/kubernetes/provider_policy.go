package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"

	"gorm.io/gorm"
)

type PolicyProvider struct {
	Repo       *repository.PolicyRepository
	Service    *service.PolicyManager
	Controller *controller.PolicyController
}

func NewPolicyProvider(db *gorm.DB) *PolicyProvider {
	// 初始化用户模块的依赖
	repo := repository.NewPolicyRepository(db)
	svc := service.NewPolicyManager(repo, db)
	ctrl := controller.NewPolicyController(svc)

	return &PolicyProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
