package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

// PipelineRunManager 定义接口
type PipelineRunManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.PipelineRun, error)
	Get(ctx context.Context, id uint, ns, name string) (*v1.PipelineRun, error)
	Create(ctx context.Context, id uint, ns string, body *v1.PipelineRun) (*v1.PipelineRun, error)
	Update(ctx context.Context, id uint, ns string, body *v1.PipelineRun) (*v1.PipelineRun, error)
	Delete(ctx context.Context, id uint, ns string, name string) error
}

type pipelineRunManager struct {
	pipelineRun repository.PipelineRunRepository
}

func NewPipelineRunManager(pipelineRun repository.PipelineRunRepository) PipelineRunManager {
	return &pipelineRunManager{pipelineRun: pipelineRun}
}

// List 列表
func (s *pipelineRunManager) List(ctx context.Context, id uint, ns string) ([]dto.PipelineRun, error) {
	items, err := s.pipelineRun.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToPipelineRunDTOs(items), nil
}

// Get 查询
func (s *pipelineRunManager) Get(ctx context.Context, id uint, ns, name string) (*v1.PipelineRun, error) {
	return s.pipelineRun.Get(ctx, id, ns, name)
}

// Create 创建
func (s *pipelineRunManager) Create(ctx context.Context, id uint, ns string, body *v1.PipelineRun) (*v1.PipelineRun, error) {
	return s.pipelineRun.Create(ctx, id, ns, body)
}

// Update 更新
func (s *pipelineRunManager) Update(ctx context.Context, id uint, ns string, body *v1.PipelineRun) (*v1.PipelineRun, error) {
	return s.pipelineRun.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *pipelineRunManager) Delete(ctx context.Context, id uint, ns, name string) (err error) {
	return s.pipelineRun.Delete(ctx, id, ns, name)
}
