package service

import (
	"context"
	"errors"
	"valyria-backend/internal/apps/dragon/dto"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/core/config"

	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

type TektonTaskRunManager interface {
	List(ctx context.Context) ([]dto.TektonTaskRunDTO, error)
	Get(ctx context.Context, name string) (*tektonv1.TaskRun, error)
	Create(ctx context.Context, body *tektonv1.TaskRun) (*tektonv1.TaskRun, error)
	Update(ctx context.Context, name string, body *tektonv1.TaskRun) (*tektonv1.TaskRun, error)
	Delete(ctx context.Context, name string) error
	BatchDelete(ctx context.Context, names []string) error
}

type tektonTaskRunManager struct {
	repo repository.TektonTaskRunRepository
	cfg  *config.Config
}

func NewTektonTaskRunManager(repo repository.TektonTaskRunRepository, cfg *config.Config) TektonTaskRunManager {
	return &tektonTaskRunManager{repo: repo, cfg: cfg}
}

func (s *tektonTaskRunManager) List(ctx context.Context) ([]dto.TektonTaskRunDTO, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	items, err := s.repo.List(ctx, ns)
	if err != nil {
		return nil, err
	}
	result := make([]dto.TektonTaskRunDTO, 0, len(items))
	for _, item := range items {
		result = append(result, dto.TektonTaskRunDTO{
			Name:      item.Name,
			Namespace: item.Namespace,
			Succeeded: item.Succeeded,
			Reason:    item.Reason,
			CreatedAt: item.CreatedAt,
		})
	}
	return result, nil
}

func (s *tektonTaskRunManager) Get(ctx context.Context, name string) (*tektonv1.TaskRun, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, ns, name)
}

func (s *tektonTaskRunManager) Create(ctx context.Context, body *tektonv1.TaskRun) (*tektonv1.TaskRun, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	body.Namespace = ns
	return s.repo.Create(ctx, ns, body)
}

func (s *tektonTaskRunManager) Update(ctx context.Context, name string, body *tektonv1.TaskRun) (*tektonv1.TaskRun, error) {
	ns, err := s.getNamespace()
	if err != nil {
		return nil, err
	}
	body.Name = name
	body.Namespace = ns
	return s.repo.Update(ctx, ns, body)
}

func (s *tektonTaskRunManager) Delete(ctx context.Context, name string) error {
	ns, err := s.getNamespace()
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, ns, name)
}

func (s *tektonTaskRunManager) BatchDelete(ctx context.Context, names []string) error {
	ns, err := s.getNamespace()
	if err != nil {
		return err
	}
	return s.repo.BatchDelete(ctx, ns, names)
}

func (s *tektonTaskRunManager) getNamespace() (string, error) {
	if s.cfg == nil || s.cfg.K8s.Tekton.Namespace == "" {
		return "", errors.New("tekton namespace is required")
	}
	return s.cfg.K8s.Tekton.Namespace, nil
}
