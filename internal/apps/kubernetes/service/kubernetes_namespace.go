package service

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"valyria-backend/internal/apps/kubernetes/repository"
)

// NamespaceService 定义接口
type NamespaceService interface {
	List(ctx context.Context, id int) ([]repository.Namespace, error)
	Create(ctx context.Context, id int, body *v1.Namespace) (namespace *v1.Namespace, err error)
}

type namespaceService struct {
	namespace repository.NamespaceRepository
}

func NewNamespaceService(namespace repository.NamespaceRepository) NamespaceService {
	return &namespaceService{namespace: namespace}
}

// List 列表
func (s *namespaceService) List(ctx context.Context, id int) (namespaces []repository.Namespace, err error) {
	return s.namespace.List(ctx, id)
}

// Create 创建命名空间
func (s *namespaceService) Create(ctx context.Context, id int, body *v1.Namespace) (namespace *v1.Namespace, err error) {
	return s.namespace.Create(ctx, id, body)
}
