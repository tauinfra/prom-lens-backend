package prometheus

import (
	"valyria-backend/internal/apps/prometheus/controller"
	"valyria-backend/internal/apps/prometheus/executor"
	"valyria-backend/internal/apps/prometheus/repository"
	"valyria-backend/internal/apps/prometheus/service"

	"gorm.io/gorm"
)

type RuleProvider struct {
	Repo       *repository.RuleRepository
	Service    *service.RuleManager
	Controller *controller.RuleController
}

func NewRuleProvider(db *gorm.DB, syncer executor.RuleSyncer) *RuleProvider {
	repo := repository.NewRuleRepository(db)
	svc := service.NewRuleManager(repo, syncer)
	ctrl := controller.NewRuleController(svc)

	return &RuleProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
