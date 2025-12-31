package service

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	"valyria-backend/internal/apps/kubernetes/repository"
)

// ReplicaSetManager 定义接口
type ReplicaSetManager interface {
	List(ctx context.Context, id int, ns, labelSelector string) ([]repository.ReplicaSet, error)
	Get(ctx context.Context, id int, ns, name string) (*appsv1.ReplicaSet, error)
	Update(ctx context.Context, id int, ns string, body *appsv1.ReplicaSet) (*appsv1.ReplicaSet, error)
	Delete(ctx context.Context, id int, ns string, name string) error
}

type replicaSetManager struct {
	replicaset repository.ReplicaSetRepository
}

func NewReplicaSetManager(replicaset repository.ReplicaSetRepository) ReplicaSetManager {
	return &replicaSetManager{replicaset: replicaset}
}

// List 列表
func (s *replicaSetManager) List(ctx context.Context, id int, ns, labelSelector string) ([]repository.ReplicaSet, error) {
	return s.replicaset.List(ctx, id, ns, labelSelector)
}

// Get 查询
func (s *replicaSetManager) Get(ctx context.Context, id int, ns, name string) (*appsv1.ReplicaSet, error) {
	return s.replicaset.Get(ctx, id, ns, name)
}

// Update 更新
func (s *replicaSetManager) Update(ctx context.Context, id int, ns string, body *appsv1.ReplicaSet) (*appsv1.ReplicaSet, error) {
	return s.replicaset.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *replicaSetManager) Delete(ctx context.Context, id int, ns, name string) (err error) {
	return s.replicaset.Delete(ctx, id, ns, name)
}
