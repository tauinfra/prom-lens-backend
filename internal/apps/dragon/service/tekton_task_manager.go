package service

import (
	"context"
	"errors"
	"valyria-backend/internal/apps/dragon/dto"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/core/config"

	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

type TektonTaskManager interface {
	List(ctx context.Context) ([]dto.TektonTaskDTO, error)
	Get(ctx context.Context, name string) (*tektonv1.Task, error)
	Create(ctx context.Context, body *tektonv1.Task) (*tektonv1.Task, error)
	Update(ctx context.Context, name string, body *tektonv1.Task) (*tektonv1.Task, error)
	Delete(ctx context.Context, name string) error
}

type tektonTaskManager struct {
	repo repository.TektonTaskRepository
	cfg  *config.Config
}

func NewTektonTaskManager(repo repository.TektonTaskRepository, cfg *config.Config) TektonTaskManager {
	return &tektonTaskManager{repo: repo, cfg: cfg}
}

func (s *tektonTaskManager) List(ctx context.Context) ([]dto.TektonTaskDTO, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	items, err := s.repo.List(ctx, ns)
	if err != nil {
		return nil, err
	}
	result := make([]dto.TektonTaskDTO, 0, len(items))
	for _, item := range items {
		result = append(result, dto.TektonTaskDTO{
			Name:      item.Name,
			Namespace: item.Namespace,
			CreatedAt: item.CreatedAt,
		})
	}
	return result, nil
}

func (s *tektonTaskManager) Get(ctx context.Context, name string) (*tektonv1.Task, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, ns, name)
}

func (s *tektonTaskManager) Create(ctx context.Context, body *tektonv1.Task) (*tektonv1.Task, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	body.Namespace = ns
	return s.repo.Create(ctx, ns, body)
}

func (s *tektonTaskManager) Update(ctx context.Context, name string, body *tektonv1.Task) (*tektonv1.Task, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	body.Name = name
	body.Namespace = ns
	return s.repo.Update(ctx, ns, body)
}

func (s *tektonTaskManager) Delete(ctx context.Context, name string) error {
	ns, err := s.getNamespace()
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, ns, name)
}

func (s *tektonTaskManager) getNamespace() (string, error) {
	if s.cfg == nil || s.cfg.K8s.Tekton.Namespace == "" {
		return "", errors.New("tekton namespace is required")
	}
	return s.cfg.K8s.Tekton.Namespace, nil
}
