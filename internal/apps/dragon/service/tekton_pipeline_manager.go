package service

import (
	"context"
	"errors"
	"valyria-backend/internal/apps/dragon/dto"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/core/config"

	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

type TektonPipelineManager interface {
	List(ctx context.Context) ([]dto.TektonPipelineDTO, error)
	Get(ctx context.Context, name string) (*tektonv1.Pipeline, error)
	Create(ctx context.Context, body *tektonv1.Pipeline) (*tektonv1.Pipeline, error)
	Update(ctx context.Context, name string, body *tektonv1.Pipeline) (*tektonv1.Pipeline, error)
	Delete(ctx context.Context, name string) error
}

type tektonPipelineManager struct {
	repo repository.TektonPipelineRepository
	cfg  *config.Config
}

func NewTektonPipelineManager(repo repository.TektonPipelineRepository, cfg *config.Config) TektonPipelineManager {
	return &tektonPipelineManager{repo: repo, cfg: cfg}
}

func (s *tektonPipelineManager) List(ctx context.Context) ([]dto.TektonPipelineDTO, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	items, err := s.repo.List(ctx, ns)
	if err != nil {
		return nil, err
	}
	result := make([]dto.TektonPipelineDTO, 0, len(items))
	for _, item := range items {
		result = append(result, dto.TektonPipelineDTO{
			Name:      item.Name,
			Namespace: item.Namespace,
			CreatedAt: item.CreatedAt,
		})
	}
	return result, nil
}

func (s *tektonPipelineManager) Get(ctx context.Context, name string) (*tektonv1.Pipeline, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, ns, name)
}

func (s *tektonPipelineManager) Create(ctx context.Context, body *tektonv1.Pipeline) (*tektonv1.Pipeline, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	body.Namespace = ns
	return s.repo.Create(ctx, ns, body)
}

func (s *tektonPipelineManager) Update(ctx context.Context, name string, body *tektonv1.Pipeline) (*tektonv1.Pipeline, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	body.Name = name
	body.Namespace = ns
	return s.repo.Update(ctx, ns, body)
}

func (s *tektonPipelineManager) Delete(ctx context.Context, name string) error {
	ns, err := s.getNamespace()
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, ns, name)
}

func (s *tektonPipelineManager) getNamespace() (string, error) {
	if s.cfg == nil || s.cfg.K8s.Tekton.Namespace == "" {
		return "", errors.New("tekton namespace is required")
	}
	return s.cfg.K8s.Tekton.Namespace, nil
}
