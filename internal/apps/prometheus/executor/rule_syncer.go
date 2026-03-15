package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"valyria-backend/internal/apps/prometheus/model"

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
	var data []byte
	//
	group, err := s.fetchGroupWithRules(ctx, groupID)
	if err != nil {
		return err
	}
	//
	switch group.Type {
	case AlertingRules:
		data, err = s.buildPromRuleGroupYAML(group.Name, group.Rules)
	default:
		return fmt.Errorf("unsupported group type: %s", group.Type)
	}
	//
	return s.syncToConfigMap(ctx, group.Name, data)
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
			rules = append(rules, Rules{
				Alert:  r.Name,
				For:    r.For,
				Expr:   r.Expr,
				Labels: r.Labels,
				Annotations: Annotations{
					Summary:     r.Summary,
					Description: r.Description,
				},
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

// fetchGroupWithRules 获取 Prom 规则组数据
func (s *configMapSyncer) fetchGroupWithRules(ctx context.Context, groupID int) (model.Group, error) {
	var group model.Group
	err := s.db.WithContext(ctx).Preload("Rules").First(&group, groupID).Error
	if err != nil {
		return model.Group{}, fmt.Errorf("failed to fetch group %d: %w", groupID, err)
	}
	return group, nil
}

// syncToConfigMap 更新 Prom configmap 配置文件
func (s *configMapSyncer) syncToConfigMap(ctx context.Context, groupName string, data []byte) error {
	client, err := clientSet()
	if err != nil {
		return err
	}
	promNamespace, promConfigMap := getPromRuleConfig()
	// 获取配置
	cm, err := client.CoreV1().ConfigMaps(promNamespace).Get(ctx, promConfigMap, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	key := fmt.Sprintf("%s.yml", groupName)
	cm.Data[key] = string(data)
	// 更新配置
	_, err = client.CoreV1().ConfigMaps(promNamespace).Update(ctx, cm, metav1.UpdateOptions{})
	return err
}

 
