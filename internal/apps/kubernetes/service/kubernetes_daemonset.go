package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/repository"

	appsv1 "k8s.io/api/apps/v1"
)

// DaemonSetManager 定义接口
type DaemonSetManager interface {
	List(ctx context.Context, id int, ns string) ([]repository.DaemonSet, error)
	Get(ctx context.Context, id int, ns, name string) (*appsv1.DaemonSet, error)
	GetDetail(ctx context.Context, id int, ns, name string) (repository.DaemonSet, error)
	Create(ctx context.Context, id int, ns string, body *appsv1.DaemonSet) (*appsv1.DaemonSet, error)
	Update(ctx context.Context, id int, ns string, body *appsv1.DaemonSet) (*appsv1.DaemonSet, error)
	Delete(ctx context.Context, id int, ns string, name string) error
}

type daemonSetManager struct {
	daemonset repository.DaemonSetRepository
}

func NewDaemonSetManager(daemonset repository.DaemonSetRepository) DaemonSetManager {
	return &daemonSetManager{daemonset: daemonset}
}

// List 列表
func (s *daemonSetManager) List(ctx context.Context, id int, ns string) ([]repository.DaemonSet, error) {
	return s.daemonset.List(ctx, id, ns)
}

// Get 查询
func (s *daemonSetManager) Get(ctx context.Context, id int, ns, name string) (*appsv1.DaemonSet, error) {
	return s.daemonset.Get(ctx, id, ns, name)
}

// GetDetail 详情
func (s *daemonSetManager) GetDetail(ctx context.Context, id int, ns, name string) (repository.DaemonSet, error) {
	return s.daemonset.GetDetail(ctx, id, ns, name)
}

// Create 创建
func (s *daemonSetManager) Create(ctx context.Context, id int, ns string, body *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return s.daemonset.Create(ctx, id, ns, body)
}

// Update 更新
func (s *daemonSetManager) Update(ctx context.Context, id int, ns string, body *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return s.daemonset.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *daemonSetManager) Delete(ctx context.Context, id int, ns, name string) (err error) {
	return s.daemonset.Delete(ctx, id, ns, name)
}
