package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	networkingv1 "k8s.io/api/networking/v1"
)

type IngressClassManager interface {
	List(ctx context.Context, id uint) ([]dto.IngressClass, error)
	Get(ctx context.Context, id uint, name string) (*networkingv1.IngressClass, error)
	Create(ctx context.Context, id uint, body *networkingv1.IngressClass) (*networkingv1.IngressClass, error)
	Update(ctx context.Context, id uint, body *networkingv1.IngressClass) (*networkingv1.IngressClass, error)
	Delete(ctx context.Context, id uint, name string) error
}

type ingressClassManager struct {
	ingressClass repository.IngressClassRepository
}

func NewIngressClassManager(ingressClass repository.IngressClassRepository) IngressClassManager {
	return &ingressClassManager{ingressClass: ingressClass}
}

func (s *ingressClassManager) List(ctx context.Context, id uint) ([]dto.IngressClass, error) {
	items, err := s.ingressClass.List(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.ToIngressClassDTOs(items), nil
}

func (s *ingressClassManager) Get(ctx context.Context, id uint, name string) (*networkingv1.IngressClass, error) {
	return s.ingressClass.Get(ctx, id, name)
}

func (s *ingressClassManager) Create(ctx context.Context, id uint, body *networkingv1.IngressClass) (*networkingv1.IngressClass, error) {
	return s.ingressClass.Create(ctx, id, body)
}

func (s *ingressClassManager) Update(ctx context.Context, id uint, body *networkingv1.IngressClass) (*networkingv1.IngressClass, error) {
	return s.ingressClass.Update(ctx, id, body)
}

func (s *ingressClassManager) Delete(ctx context.Context, id uint, name string) error {
	return s.ingressClass.Delete(ctx, id, name)
}
