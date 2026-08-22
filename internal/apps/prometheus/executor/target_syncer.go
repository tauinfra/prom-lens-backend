package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"prom-lens-backend/internal/apps/prometheus/model"
	"prom-lens-backend/internal/core/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

type fileSDTarget struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func (s *configMapTargetSyncer) buildTargetGroupJSON(group model.TargetGroup, targets []model.Target) ([]byte, error) {
	entries := make([]fileSDTarget, 0, len(targets))
	for _, item := range targets {
		if item.Enabled != nil && !*item.Enabled {
			continue
		}
		labels, err := mergeTargetLabels(group.Labels, item.Labels)
		if err != nil {
			return nil, fmt.Errorf("merge labels for %s:%d: %w", item.IPAddress, item.Port, err)
		}
		entries = append(entries, fileSDTarget{
			Targets: []string{fmt.Sprintf("%s:%d", item.IPAddress, item.Port)},
			Labels:  labels,
		})
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal file_sd targets failed: %w", err)
	}
	return data, nil
}

func mergeTargetLabels(groupLabels, targetLabels datatypes.JSON) (map[string]string, error) {
	out := make(map[string]string)
	if len(groupLabels) > 0 {
		if err := json.Unmarshal(groupLabels, &out); err != nil {
			return nil, err
		}
	}
	if len(targetLabels) > 0 {
		var target map[string]string
		if err := json.Unmarshal(targetLabels, &target); err != nil {
			return nil, err
		}
		for k, v := range target {
			out[k] = v
		}
	}
	if out == nil {
		out = map[string]string{}
	}
	return out, nil
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
	if errors.IsNotFound(err) {
		cm = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      promConfigMap,
				Namespace: promNamespace,
			},
			Data: map[string]string{key: string(data)},
		}
		if _, err = client.CoreV1().ConfigMaps(promNamespace).Create(ctx, cm, metav1.CreateOptions{}); err != nil {
			logPromSyncFailed("target", groupID, groupName, "create_configmap", err)
			return fmt.Errorf("create target configmap %s/%s: %w", promNamespace, promConfigMap, err)
		}
		logPromSyncSuccess("target", promNamespace, promConfigMap, key, len(data))
		return nil
	}
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
