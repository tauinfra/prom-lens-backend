package prometheus

import (
	"prom-lens-backend/internal/apps/prometheus/controller"
	"prom-lens-backend/internal/apps/prometheus/executor"
	"prom-lens-backend/internal/apps/prometheus/repository"
	"prom-lens-backend/internal/apps/prometheus/service"

	"gorm.io/gorm"
)

type RuleProvider struct {
	Repo       *repository.RuleRepository
	Service    *service.RuleManager
	Controller *controller.RuleController
}

func NewRuleProvider(db *gorm.DB, syncer executor.RuleSyncer) *RuleProvider {
	groupRepo := repository.NewGroupRepository(db)
	repo := repository.NewRuleRepository(db)
	svc := service.NewRuleManager(groupRepo, repo, syncer)
	ctrl := controller.NewRuleController(svc)

	return &RuleProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
