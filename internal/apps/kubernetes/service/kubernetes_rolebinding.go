package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	rbacv1 "k8s.io/api/rbac/v1"
)

type RoleBindingManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.RoleBinding, error)
	Get(ctx context.Context, id uint, ns, name string) (*rbacv1.RoleBinding, error)
	Create(ctx context.Context, id uint, ns string, body *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error)
	Update(ctx context.Context, id uint, ns string, body *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type roleBindingManager struct {
	roleBinding repository.RoleBindingRepository
}

func NewRoleBindingManager(roleBinding repository.RoleBindingRepository) RoleBindingManager {
	return &roleBindingManager{roleBinding: roleBinding}
}

func (s *roleBindingManager) List(ctx context.Context, id uint, ns string) ([]dto.RoleBinding, error) {
	items, err := s.roleBinding.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToRoleBindingDTOs(items), nil
}

func (s *roleBindingManager) Get(ctx context.Context, id uint, ns, name string) (*rbacv1.RoleBinding, error) {
	return s.roleBinding.Get(ctx, id, ns, name)
}

func (s *roleBindingManager) Create(ctx context.Context, id uint, ns string, body *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error) {
	return s.roleBinding.Create(ctx, id, ns, body)
}

func (s *roleBindingManager) Update(ctx context.Context, id uint, ns string, body *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error) {
	return s.roleBinding.Update(ctx, id, ns, body)
}

func (s *roleBindingManager) Delete(ctx context.Context, id uint, ns, name string) error {
	return s.roleBinding.Delete(ctx, id, ns, name)
}
