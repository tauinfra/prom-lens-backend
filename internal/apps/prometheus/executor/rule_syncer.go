package executor

import (
	"context"
	"encoding/json"
	"fmt"

	prom "prom-lens-backend/internal/apps/prometheus"
	"prom-lens-backend/internal/apps/prometheus/model"

	"prom-lens-backend/internal/core/logger"

	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// RuleSyncer 规则同步接口
type RuleSyncer interface {
	SyncRuleGroup(ctx context.Context, groupID int) error
}

// configMapSyncer 默认同步实现
type configMapSyncer struct {
	db *gorm.DB
}

// NewRuleSyncer 创建默认同步实现
func NewRuleSyncer(db *gorm.DB) RuleSyncer {
	return &configMapSyncer{db: db}
}

func (s *configMapSyncer) SyncRuleGroup(ctx context.Context, groupID int) error {
	group, err := s.fetchGroup(ctx, groupID)
	if err != nil {
		logPromSyncFailed("rule", groupID, "", "fetch_group", err)
		return err
	}

	var data []byte
	switch group.Type {
	case AlertingRules:
		logger.Infof(
			"[prom-sync] start resource=rule groupID=%d groupName=%q type=%s ruleCount=%d",
			groupID, group.Name, group.Type, len(group.Rules),
		)
		data, err = s.buildPromRuleGroupYAML(group.Name, group.Rules)
		if err != nil {
			logPromSyncFailed("rule", groupID, group.Name, "build_yaml", err)
			return err
		}
	case AlertingRecords:
		logger.Infof(
			"[prom-sync] start resource=record groupID=%d groupName=%q type=%s recordCount=%d",
			groupID, group.Name, group.Type, len(group.Records),
		)
		data, err = s.buildPromRecordGroupYAML(group.Name, group.Records)
		if err != nil {
			logPromSyncFailed("record", groupID, group.Name, "build_yaml", err)
			return err
		}
	default:
		err = fmt.Errorf("unsupported group type: %s", group.Type)
		logPromSyncFailed("rule", groupID, group.Name, "validate_group_type", err)
		return err
	}

	return s.syncRuleToConfigMap(ctx, groupID, group.Name, data)
}

// marshalPromRuleToYAML 转换规则组为 YAML 格式
func (s *configMapSyncer) marshalPromRuleToYAML(data any) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal prom rule failed: %w", err)
	}

	yamlData, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return nil, fmt.Errorf("convert prom rule to yaml failed: %w", err)
	}

	return yamlData, nil
}

// buildPromRuleGroupYAML 生成规则组
func (s *configMapSyncer) buildPromRuleGroupYAML(groupName string, data []model.Rule) ([]byte, error) {
	var rules []Rules

	for _, r := range data {
		if r.Status == nil || *r.Status {
			annotations, err := prom.MergeRuleAnnotations(r.Summary, r.Description, r.ExtraAnnotations)
			if err != nil {
				return nil, fmt.Errorf("rule %q annotations: %w", r.Name, err)
			}
			rules = append(rules, Rules{
				Alert:       r.Name,
				For:         r.For,
				Expr:        r.Expr,
				Labels:      r.Labels,
				Annotations: annotations,
			})
		}
	}

	ruleConfig := RuleConfig{
		Groups: []RuleGroups{
			{
				Name:  groupName,
				Rules: rules,
			},
		},
	}
	return s.marshalPromRuleToYAML(ruleConfig)
}

func (s *configMapSyncer) buildPromRecordGroupYAML(groupName string, data []model.Record) ([]byte, error) {
	records := make([]Records, 0, len(data))
	for _, item := range data {
		records = append(records, Records{
			Record: item.Name,
			Expr:   item.Expr,
		})
	}
	cfg := RecordConfig{
		Groups: []RecordGroups{
			{
				Name:  groupName,
				Rules: records,
			},
		},
	}
	return s.marshalPromRuleToYAML(cfg)
}

func (s *configMapSyncer) fetchGroup(ctx context.Context, groupID int) (model.Group, error) {
	var group model.Group
	if err := s.db.WithContext(ctx).First(&group, groupID).Error; err != nil {
		return model.Group{}, fmt.Errorf("failed to fetch group %d: %w", groupID, err)
	}
	switch group.Type {
	case AlertingRecords:
		if err := s.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&group.Records).Error; err != nil {
			return model.Group{}, fmt.Errorf("failed to fetch records for group %d: %w", groupID, err)
		}
	default:
		if err := s.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&group.Rules).Error; err != nil {
			return model.Group{}, fmt.Errorf("failed to fetch rules for group %d: %w", groupID, err)
		}
	}
	return group, nil
}

// syncRuleToConfigMap 更新 Prom 规则 ConfigMap。
func (s *configMapSyncer) syncRuleToConfigMap(ctx context.Context, groupID int, groupName string, data []byte) error {
	client, err := clientSet()
	if err != nil {
		logPromSyncFailed("rule", groupID, groupName, "k8s_client", err)
		return err
	}
	promNamespace, promConfigMap := getPromRuleConfig()
	key := fmt.Sprintf("%s.yml", groupName)

	cm, err := client.CoreV1().ConfigMaps(promNamespace).Get(ctx, promConfigMap, metav1.GetOptions{})
	if err != nil {
		logPromSyncFailed("rule", groupID, groupName, "get_configmap", err)
		return err
	}

	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	cm.Data[key] = string(data)
	if _, err = client.CoreV1().ConfigMaps(promNamespace).Update(ctx, cm, metav1.UpdateOptions{}); err != nil {
		logPromSyncFailed("rule", groupID, groupName, "update_configmap", err)
		return err
	}

	logPromSyncSuccess("rule", promNamespace, promConfigMap, key, len(data))
	return nil
}
