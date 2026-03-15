package service

import (
	"context"
	"errors"
	"strings"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/request"

	autoscalingv2 "k8s.io/api/autoscaling/v2"
)

type HPAManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.HPA, error)
	Get(ctx context.Context, id uint, ns, name string) (*autoscalingv2.HorizontalPodAutoscaler, error)
	Create(ctx context.Context, id uint, ns string, hpa *autoscalingv2.HorizontalPodAutoscaler) (*autoscalingv2.HorizontalPodAutoscaler, error)
	Update(ctx context.Context, id uint, ns string, hpa *autoscalingv2.HorizontalPodAutoscaler) (*autoscalingv2.HorizontalPodAutoscaler, error)
	Delete(ctx context.Context, id uint, ns, name string) error
	DeleteBatch(ctx context.Context, id uint, ns string, req *request.BatchDeleteWorkloadRequest) error
}

type hpaManager struct {
	repo repository.HpaRepository
}

func NewHPAManager(repo repository.HpaRepository) HPAManager {
	return &hpaManager{repo: repo}
}

func (s *hpaManager) List(ctx context.Context, id uint, ns string) ([]dto.HPA, error) {
	items, err := s.repo.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToHpaDTOs(items), nil
}

func (s *hpaManager) Get(ctx context.Context, id uint, ns, name string) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	return s.repo.Get(ctx, id, ns, name)
}

func (s *hpaManager) Create(ctx context.Context, id uint, ns string, hpa *autoscalingv2.HorizontalPodAutoscaler) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	return s.repo.Create(ctx, id, ns, hpa)
}

func (s *hpaManager) Update(ctx context.Context, id uint, ns string, hpa *autoscalingv2.HorizontalPodAutoscaler) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	return s.repo.Update(ctx, id, ns, hpa)
}

func (s *hpaManager) Delete(ctx context.Context, id uint, ns, name string) error {
	return s.repo.Delete(ctx, id, ns, name)
}

func (s *hpaManager) DeleteBatch(ctx context.Context, id uint, ns string, req *request.BatchDeleteWorkloadRequest) error {
	if req == nil || len(req.Names) == 0 {
		return errors.New("hpa names is required")
	}
	seen := make(map[string]struct{}, len(req.Names))
	for _, name := range req.Names {
		name = strings.TrimSpace(name)
		if name == "" {
			return errors.New("hpa name is required")
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		if err := s.repo.Delete(ctx, id, ns, name); err != nil {
			return err
		}
	}
	return nil
}
