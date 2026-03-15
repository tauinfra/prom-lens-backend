package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	rbacv1 "k8s.io/api/rbac/v1"
)

type ClusterRoleBindingManager interface {
	List(ctx context.Context, id uint) ([]dto.ClusterRoleBinding, error)
	Get(ctx context.Context, id uint, name string) (*rbacv1.ClusterRoleBinding, error)
	Create(ctx context.Context, id uint, body *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error)
	Update(ctx context.Context, id uint, body *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error)
	Delete(ctx context.Context, id uint, name string) error
}

type clusterRoleBindingManager struct {
	clusterRoleBinding repository.ClusterRoleBindingRepository
}

func NewClusterRoleBindingManager(clusterRoleBinding repository.ClusterRoleBindingRepository) ClusterRoleBindingManager {
	return &clusterRoleBindingManager{clusterRoleBinding: clusterRoleBinding}
}

func (s *clusterRoleBindingManager) List(ctx context.Context, id uint) ([]dto.ClusterRoleBinding, error) {
	items, err := s.clusterRoleBinding.List(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.ToClusterRoleBindingDTOs(items), nil
}

func (s *clusterRoleBindingManager) Get(ctx context.Context, id uint, name string) (*rbacv1.ClusterRoleBinding, error) {
	return s.clusterRoleBinding.Get(ctx, id, name)
}

func (s *clusterRoleBindingManager) Create(ctx context.Context, id uint, body *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error) {
	return s.clusterRoleBinding.Create(ctx, id, body)
}

func (s *clusterRoleBindingManager) Update(ctx context.Context, id uint, body *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error) {
	return s.clusterRoleBinding.Update(ctx, id, body)
}

func (s *clusterRoleBindingManager) Delete(ctx context.Context, id uint, name string) error {
	return s.clusterRoleBinding.Delete(ctx, id, name)
}
