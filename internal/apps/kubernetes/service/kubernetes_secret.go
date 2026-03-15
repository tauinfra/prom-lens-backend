package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	corev1 "k8s.io/api/core/v1"
)

// SecretManager 定义接口
type SecretManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.Secret, error)
	Get(ctx context.Context, id uint, ns, name string) (*corev1.Secret, error)
	Create(ctx context.Context, id uint, ns string, body *corev1.Secret) (*corev1.Secret, error)
	Update(ctx context.Context, id uint, ns string, body *corev1.Secret) (*corev1.Secret, error)
	Delete(ctx context.Context, id uint, ns string, name string) error
}

type secretManager struct {
	secret repository.SecretRepository
}

func NewSecretManager(secret repository.SecretRepository) SecretManager {
	return &secretManager{secret: secret}
}

// List 列表
func (s *secretManager) List(ctx context.Context, id uint, ns string) ([]dto.Secret, error) {
	items, err := s.secret.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToSecretDTOs(items), nil
}

// Get 查询
func (s *secretManager) Get(ctx context.Context, id uint, ns, name string) (*corev1.Secret, error) {
	return s.secret.Get(ctx, id, ns, name)
}

// Create 创建
func (s *secretManager) Create(ctx context.Context, id uint, ns string, body *corev1.Secret) (*corev1.Secret, error) {
	return s.secret.Create(ctx, id, ns, body)
}

// Update 更新
func (s *secretManager) Update(ctx context.Context, id uint, ns string, body *corev1.Secret) (*corev1.Secret, error) {
	return s.secret.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *secretManager) Delete(ctx context.Context, id uint, ns, name string) (err error) {
	return s.secret.Delete(ctx, id, ns, name)
}
