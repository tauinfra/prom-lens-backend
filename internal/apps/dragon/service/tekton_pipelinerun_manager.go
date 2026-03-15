package service

import (
	"context"
	"errors"
	"valyria-backend/internal/apps/dragon/dto"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/core/config"

	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

type TektonPipelineRunManager interface {
	List(ctx context.Context) ([]dto.TektonPipelineRunDTO, error)
	Get(ctx context.Context, name string) (*tektonv1.PipelineRun, error)
	Create(ctx context.Context, body *tektonv1.PipelineRun) (*tektonv1.PipelineRun, error)
	Update(ctx context.Context, name string, body *tektonv1.PipelineRun) (*tektonv1.PipelineRun, error)
	Delete(ctx context.Context, name string) error
	BatchDelete(ctx context.Context, names []string) error
}

type tektonPipelineRunManager struct {
	repo repository.TektonPipelineRunRepository
	cfg  *config.Config
}

func NewTektonPipelineRunManager(repo repository.TektonPipelineRunRepository, cfg *config.Config) TektonPipelineRunManager {
	return &tektonPipelineRunManager{repo: repo, cfg: cfg}
}

func (s *tektonPipelineRunManager) List(ctx context.Context) ([]dto.TektonPipelineRunDTO, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	items, err := s.repo.List(ctx, ns)
	if err != nil {
		return nil, err
	}
	result := make([]dto.TektonPipelineRunDTO, 0, len(items))
	for _, item := range items {
		result = append(result, dto.TektonPipelineRunDTO{
			Name:      item.Name,
			Namespace: item.Namespace,
			Succeeded: item.Succeeded,
			Reason:    item.Reason,
			CreatedAt: item.CreatedAt,
		})
	}
	return result, nil
}

func (s *tektonPipelineRunManager) Get(ctx context.Context, name string) (*tektonv1.PipelineRun, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, ns, name)
}

func (s *tektonPipelineRunManager) Create(ctx context.Context, body *tektonv1.PipelineRun) (*tektonv1.PipelineRun, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	body.Namespace = ns
	return s.repo.Create(ctx, ns, body)
}

func (s *tektonPipelineRunManager) Update(ctx context.Context, name string, body *tektonv1.PipelineRun) (*tektonv1.PipelineRun, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	body.Name = name
	body.Namespace = ns
	return s.repo.Update(ctx, ns, body)
}

func (s *tektonPipelineRunManager) Delete(ctx context.Context, name string) error {
	ns, err := s.getNamespace()
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, ns, name)
}

func (s *tektonPipelineRunManager) BatchDelete(ctx context.Context, names []string) error {
	ns, err := s.getNamespace()
	if err != nil {
		return err
	}
	return s.repo.BatchDelete(ctx, ns, names)
}

func (s *tektonPipelineRunManager) getNamespace() (string, error) {
	if s.cfg == nil || s.cfg.K8s.Tekton.Namespace == "" {
		return "", errors.New("tekton namespace is required")
	}
	return s.cfg.K8s.Tekton.Namespace, nil
}
