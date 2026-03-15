package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	rbacv1 "k8s.io/api/rbac/v1"
)

type ClusterRoleManager interface {
	List(ctx context.Context, id uint) ([]dto.ClusterRole, error)
	Get(ctx context.Context, id uint, name string) (*rbacv1.ClusterRole, error)
	Create(ctx context.Context, id uint, body *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error)
	Update(ctx context.Context, id uint, body *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error)
	Delete(ctx context.Context, id uint, name string) error
}

type clusterRoleManager struct {
	role repository.ClusterRoleRepository
}

func NewClusterRoleManager(role repository.ClusterRoleRepository) ClusterRoleManager {
	return &clusterRoleManager{role: role}
}

func (s *clusterRoleManager) List(ctx context.Context, id uint) ([]dto.ClusterRole, error) {
	items, err := s.role.List(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.ToClusterRoleDTOs(items), nil
}

func (s *clusterRoleManager) Get(ctx context.Context, id uint, name string) (*rbacv1.ClusterRole, error) {
	return s.role.Get(ctx, id, name)
}

func (s *clusterRoleManager) Create(ctx context.Context, id uint, body *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error) {
	return s.role.Create(ctx, id, body)
}

func (s *clusterRoleManager) Update(ctx context.Context, id uint, body *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error) {
	return s.role.Update(ctx, id, body)
}

func (s *clusterRoleManager) Delete(ctx context.Context, id uint, name string) error {
	return s.role.Delete(ctx, id, name)
}
