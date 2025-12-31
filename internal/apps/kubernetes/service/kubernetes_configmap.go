package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/repository"

	corev1 "k8s.io/api/core/v1"
)

// ConfigmapService 定义接口
type ConfigmapService interface {
	List(ctx context.Context, id int, ns string) ([]repository.ConfigMap, error)
	Get(ctx context.Context, id int, ns, name string) (*corev1.ConfigMap, error)
	Create(ctx context.Context, id int, ns string, body *corev1.ConfigMap) (*corev1.ConfigMap, error)
	Update(ctx context.Context, id int, ns string, body *corev1.ConfigMap) (*corev1.ConfigMap, error)
	Delete(ctx context.Context, id int, ns string, name string) error
}

type configmapService struct {
	configmap repository.ConfigmapRepository
}

func NewConfigmapService(configmap repository.ConfigmapRepository) ConfigmapService {
	return &configmapService{configmap: configmap}
}

// List 列表
func (s *configmapService) List(ctx context.Context, id int, ns string) ([]repository.ConfigMap, error) {
	return s.configmap.List(ctx, id, ns)
}

// Get 查询
func (s *configmapService) Get(ctx context.Context, id int, ns, name string) (*corev1.ConfigMap, error) {
	return s.configmap.Get(ctx, id, ns, name)
}

// Create 创建
func (s *configmapService) Create(ctx context.Context, id int, ns string, body *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	return s.configmap.Create(ctx, id, ns, body)
}

// Update 更新
func (s *configmapService) Update(ctx context.Context, id int, ns string, body *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	return s.configmap.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *configmapService) Delete(ctx context.Context, id int, ns, name string) (err error) {
	return s.configmap.Delete(ctx, id, ns, name)
}
