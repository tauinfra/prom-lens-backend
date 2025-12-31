package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/repository"

	corev1 "k8s.io/api/core/v1"
)

// ServiceManager 定义接口
type ServiceManager interface {
	List(ctx context.Context, id int, ns, labelSelector string) ([]repository.Service, error)
	Get(ctx context.Context, id int, ns, name string) (*corev1.Service, error)
	Create(ctx context.Context, id int, ns string, body *corev1.Service) (*corev1.Service, error)
	Update(ctx context.Context, id int, ns string, body *corev1.Service) (*corev1.Service, error)
	Delete(ctx context.Context, id int, ns string, name string) error
}

type serviceManager struct {
	service repository.ServiceRepository
}

func NewServiceManager(service repository.ServiceRepository) ServiceManager {
	return &serviceManager{service: service}
}

// List 列表
func (s *serviceManager) List(ctx context.Context, id int, ns, labelSelector string) ([]repository.Service, error) {
	return s.service.List(ctx, id, ns, labelSelector)
}

// Get 查询
func (s *serviceManager) Get(ctx context.Context, id int, ns, name string) (*corev1.Service, error) {
	return s.service.Get(ctx, id, ns, name)
}

// Create 创建
func (s *serviceManager) Create(ctx context.Context, id int, ns string, body *corev1.Service) (*corev1.Service, error) {
	return s.service.Create(ctx, id, ns, body)
}

// Update 更新
func (s *serviceManager) Update(ctx context.Context, id int, ns string, body *corev1.Service) (*corev1.Service, error) {
	return s.service.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *serviceManager) Delete(ctx context.Context, id int, ns, name string) (err error) {
	return s.service.Delete(ctx, id, ns, name)
}
