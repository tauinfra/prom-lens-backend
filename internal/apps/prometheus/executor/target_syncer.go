package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"prom-lens-backend/internal/apps/prometheus/model"
	"prom-lens-backend/internal/core/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TargetSyncer Target 同步接口
type TargetSyncer interface {
	SyncTargetGroup(ctx context.Context, groupID int) error
}

// configMapTargetSyncer 默认同步实现
type configMapTargetSyncer struct {
	db *gorm.DB
}

// NewTargetSyncer 创建默认同步实现
func NewTargetSyncer(db *gorm.DB) TargetSyncer {
	return &configMapTargetSyncer{db: db}
}

func (s *configMapTargetSyncer) SyncTargetGroup(ctx context.Context, groupID int) error {
	group, targets, err := s.fetchGroupWithTargets(ctx, groupID)
	if err != nil {
		logPromSyncFailed("target", groupID, "", "fetch_group", err)
		return err
	}

	logger.Infof(
		"[prom-sync] start resource=target groupID=%d groupName=%q targetCount=%d",
		groupID, group.Name, len(targets),
	)

	data, err := s.buildTargetGroupJSON(group, targets)
	if err != nil {
		logPromSyncFailed("target", groupID, group.Name, "build_json", err)
		return err
	}
	return s.syncTargetToConfigMap(ctx, groupID, group.Name, data)
}

type targetConfig struct {
	IPAddress string         `json:"ipAddress"`
	Port      int            `json:"port"`
	Labels    datatypes.JSON `json:"labels"`
}

type targetGroupConfig struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Labels      datatypes.JSON `json:"labels"`
	Targets     []targetConfig `json:"targets"`
}

func (s *configMapTargetSyncer) buildTargetGroupJSON(group model.TargetGroup, targets []model.Target) ([]byte, error) {
	result := targetGroupConfig{
		Name:        group.Name,
		Description: group.Description,
		Labels:      group.Labels,
		Targets:     make([]targetConfig, 0, len(targets)),
	}
	for _, item := range targets {
		if item.Enabled != nil && !*item.Enabled {
			continue
		}
		result.Targets = append(result.Targets, targetConfig{
			IPAddress: item.IPAddress,
			Port:      item.Port,
			Labels:    item.Labels,
		})
	}
	data, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("marshal target group failed: %w", err)
	}
	return data, nil
}

func (s *configMapTargetSyncer) fetchGroupWithTargets(ctx context.Context, groupID int) (model.TargetGroup, []model.Target, error) {
	var (
		group   model.TargetGroup
		targets []model.Target
	)
	if err := s.db.WithContext(ctx).First(&group, groupID).Error; err != nil {
		return model.TargetGroup{}, nil, fmt.Errorf("failed to fetch target group %d: %w", groupID, err)
	}
	if err := s.db.WithContext(ctx).Where("group_id = ?", groupID).Order("id ASC").Find(&targets).Error; err != nil {
		return model.TargetGroup{}, nil, fmt.Errorf("failed to fetch targets for group %d: %w", groupID, err)
	}
	return group, targets, nil
}

func (s *configMapTargetSyncer) syncTargetToConfigMap(ctx context.Context, groupID int, groupName string, data []byte) error {
	client, err := clientSet()
	if err != nil {
		logPromSyncFailed("target", groupID, groupName, "k8s_client", err)
		return err
	}
	promNamespace, promConfigMap := getPromTargetConfig()
	key := fmt.Sprintf("%s.json", groupName)

	cm, err := client.CoreV1().ConfigMaps(promNamespace).Get(ctx, promConfigMap, metav1.GetOptions{})
	if err != nil {
		logPromSyncFailed("target", groupID, groupName, "get_configmap", err)
		return err
	}

	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	cm.Data[key] = string(data)
	if _, err = client.CoreV1().ConfigMaps(promNamespace).Update(ctx, cm, metav1.UpdateOptions{}); err != nil {
		logPromSyncFailed("target", groupID, groupName, "update_configmap", err)
		return err
	}

	logPromSyncSuccess("target", promNamespace, promConfigMap, key, len(data))
	return nil
}
