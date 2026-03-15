package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	corev1 "k8s.io/api/core/v1"
)

// ConfigmapManager 定义接口
type ConfigmapManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.ConfigMap, error)
	Get(ctx context.Context, id uint, ns, name string) (*corev1.ConfigMap, error)
	Create(ctx context.Context, id uint, ns string, body *corev1.ConfigMap) (*corev1.ConfigMap, error)
	Update(ctx context.Context, id uint, ns string, body *corev1.ConfigMap) (*corev1.ConfigMap, error)
	Delete(ctx context.Context, id uint, ns string, name string) error
}

type configmapManager struct {
	configmap repository.ConfigmapRepository
}

func NewConfigmapManager(configmap repository.ConfigmapRepository) ConfigmapManager {
	return &configmapManager{configmap: configmap}
}

// List 列表
func (s *configmapManager) List(ctx context.Context, id uint, ns string) ([]dto.ConfigMap, error) {
	items, err := s.configmap.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToConfigMapDTOs(items), nil
}

// Get 查询
func (s *configmapManager) Get(ctx context.Context, id uint, ns, name string) (*corev1.ConfigMap, error) {
	return s.configmap.Get(ctx, id, ns, name)
}

// Create 创建
func (s *configmapManager) Create(ctx context.Context, id uint, ns string, body *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	return s.configmap.Create(ctx, id, ns, body)
}

// Update 更新
func (s *configmapManager) Update(ctx context.Context, id uint, ns string, body *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	return s.configmap.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *configmapManager) Delete(ctx context.Context, id uint, ns, name string) (err error) {
	return s.configmap.Delete(ctx, id, ns, name)
}
