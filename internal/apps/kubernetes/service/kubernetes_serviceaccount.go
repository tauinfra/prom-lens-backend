package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	corev1 "k8s.io/api/core/v1"
)

type ServiceAccountManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.ServiceAccount, error)
	Get(ctx context.Context, id uint, ns, name string) (*corev1.ServiceAccount, error)
	Create(ctx context.Context, id uint, ns string, body *corev1.ServiceAccount) (*corev1.ServiceAccount, error)
	Update(ctx context.Context, id uint, ns string, body *corev1.ServiceAccount) (*corev1.ServiceAccount, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type serviceAccountManager struct {
	serviceAccount repository.ServiceAccountRepository
}

func NewServiceAccountManager(serviceAccount repository.ServiceAccountRepository) ServiceAccountManager {
	return &serviceAccountManager{serviceAccount: serviceAccount}
}

func (s *serviceAccountManager) List(ctx context.Context, id uint, ns string) ([]dto.ServiceAccount, error) {
	items, err := s.serviceAccount.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToServiceAccountDTOs(items), nil
}

func (s *serviceAccountManager) Get(ctx context.Context, id uint, ns, name string) (*corev1.ServiceAccount, error) {
	return s.serviceAccount.Get(ctx, id, ns, name)
}

func (s *serviceAccountManager) Create(ctx context.Context, id uint, ns string, body *corev1.ServiceAccount) (*corev1.ServiceAccount, error) {
	return s.serviceAccount.Create(ctx, id, ns, body)
}

func (s *serviceAccountManager) Update(ctx context.Context, id uint, ns string, body *corev1.ServiceAccount) (*corev1.ServiceAccount, error) {
	return s.serviceAccount.Update(ctx, id, ns, body)
}

func (s *serviceAccountManager) Delete(ctx context.Context, id uint, ns, name string) error {
	return s.serviceAccount.Delete(ctx, id, ns, name)
}
