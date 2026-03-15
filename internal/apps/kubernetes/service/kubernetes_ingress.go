package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	networkingv1 "k8s.io/api/networking/v1"
)

type IngressManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.Ingress, error)
	Get(ctx context.Context, id uint, ns, name string) (*networkingv1.Ingress, error)
	GetDetail(ctx context.Context, id uint, ns, name string) (dto.IngressDetail, error)
	Create(ctx context.Context, id uint, ns string, body *networkingv1.Ingress) (*networkingv1.Ingress, error)
	Update(ctx context.Context, id uint, ns string, body *networkingv1.Ingress) (*networkingv1.Ingress, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type ingressManager struct {
	ingress repository.IngressRepository
}

func NewIngressManager(ingress repository.IngressRepository) IngressManager {
	return &ingressManager{ingress: ingress}
}

func (s *ingressManager) List(ctx context.Context, id uint, ns string) ([]dto.Ingress, error) {
	items, err := s.ingress.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToIngressDTOs(items), nil
}

func (s *ingressManager) Get(ctx context.Context, id uint, ns, name string) (*networkingv1.Ingress, error) {
	return s.ingress.Get(ctx, id, ns, name)
}

func (s *ingressManager) GetDetail(ctx context.Context, id uint, ns, name string) (dto.IngressDetail, error) {
	item, err := s.ingress.Get(ctx, id, ns, name)
	if err != nil {
		return dto.IngressDetail{}, err
	}
	return dto.ToIngressDetailDTO(item), nil
}

func (s *ingressManager) Create(ctx context.Context, id uint, ns string, body *networkingv1.Ingress) (*networkingv1.Ingress, error) {
	return s.ingress.Create(ctx, id, ns, body)
}

func (s *ingressManager) Update(ctx context.Context, id uint, ns string, body *networkingv1.Ingress) (*networkingv1.Ingress, error) {
	return s.ingress.Update(ctx, id, ns, body)
}

func (s *ingressManager) Delete(ctx context.Context, id uint, ns, name string) error {
	return s.ingress.Delete(ctx, id, ns, name)
}
