package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	rbacv1 "k8s.io/api/rbac/v1"
)

type RoleManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.Role, error)
	Get(ctx context.Context, id uint, ns, name string) (*rbacv1.Role, error)
	Create(ctx context.Context, id uint, ns string, body *rbacv1.Role) (*rbacv1.Role, error)
	Update(ctx context.Context, id uint, ns string, body *rbacv1.Role) (*rbacv1.Role, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type roleManager struct {
	role repository.RoleRepository
}

func NewRoleManager(role repository.RoleRepository) RoleManager {
	return &roleManager{role: role}
}

func (s *roleManager) List(ctx context.Context, id uint, ns string) ([]dto.Role, error) {
	items, err := s.role.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToRoleDTOs(items), nil
}

func (s *roleManager) Get(ctx context.Context, id uint, ns, name string) (*rbacv1.Role, error) {
	return s.role.Get(ctx, id, ns, name)
}

func (s *roleManager) Create(ctx context.Context, id uint, ns string, body *rbacv1.Role) (*rbacv1.Role, error) {
	return s.role.Create(ctx, id, ns, body)
}

func (s *roleManager) Update(ctx context.Context, id uint, ns string, body *rbacv1.Role) (*rbacv1.Role, error) {
	return s.role.Update(ctx, id, ns, body)
}

func (s *roleManager) Delete(ctx context.Context, id uint, ns, name string) error {
	return s.role.Delete(ctx, id, ns, name)
}
