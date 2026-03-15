package prometheus

import (
	"valyria-backend/internal/apps/prometheus/executor"

	"gorm.io/gorm"
)

type Provider struct {
	Group       *GroupProvider
	Record      *RecordProvider
	Rule        *RuleProvider
	TargetGroup *TargetGroupProvider
	Target      *TargetProvider
}

func NewPrometheusProvider(db *gorm.DB) *Provider {
	syncer := executor.NewRuleSyncer(db) // Prom RuleSyncer
	targetSyncer := executor.NewTargetSyncer(db)
	return &Provider{
		Group:       NewGroupProvider(db),
		Record:      NewRecordProvider(db, syncer),
		Rule:        NewRuleProvider(db, syncer),
		TargetGroup: NewTargetGroupProvider(db),
		Target:      NewTargetProvider(db, targetSyncer),
	}
}
