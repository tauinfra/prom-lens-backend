package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	v1 "k8s.io/api/core/v1"
)

// NamespaceManager 定义接口
type NamespaceManager interface {
	List(ctx context.Context, id uint) ([]dto.Namespace, error)
	Get(ctx context.Context, id uint, name string) (*v1.Namespace, error)
	Create(ctx context.Context, id uint, body *v1.Namespace) (namespace *v1.Namespace, err error)
	Update(ctx context.Context, id uint, body *v1.Namespace) (namespace *v1.Namespace, err error)
	Delete(ctx context.Context, id uint, name string) error
}

type namespaceManager struct {
	namespace repository.NamespaceRepository
}

func NewNamespaceManager(namespace repository.NamespaceRepository) NamespaceManager {
	return &namespaceManager{namespace: namespace}
}

// List 列表
func (s *namespaceManager) List(ctx context.Context, id uint) (namespaces []dto.Namespace, err error) {
	items, err := s.namespace.List(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.ToNamespaceDTOs(items), nil
}

// Create 创建命名空间
func (s *namespaceManager) Create(ctx context.Context, id uint, body *v1.Namespace) (namespace *v1.Namespace, err error) {
	// Namespace 保护策略
	body.Labels = map[string]string{
		"platform.io/protected": "true",
	}
	return s.namespace.Create(ctx, id, body)
}

func (s *namespaceManager) Get(ctx context.Context, id uint, name string) (*v1.Namespace, error) {
	return s.namespace.Get(ctx, id, name)
}

func (s *namespaceManager) Update(ctx context.Context, id uint, body *v1.Namespace) (namespace *v1.Namespace, err error) {
	return s.namespace.Update(ctx, id, body)
}

func (s *namespaceManager) Delete(ctx context.Context, id uint, name string) error {
	return s.namespace.Delete(ctx, id, name)
}
