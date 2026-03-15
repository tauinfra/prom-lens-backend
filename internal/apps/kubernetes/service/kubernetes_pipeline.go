package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

// PipelineManager 定义接口
type PipelineManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.Pipeline, error)
	Get(ctx context.Context, id uint, ns, name string) (*v1.Pipeline, error)
	Create(ctx context.Context, id uint, ns string, body *v1.Pipeline) (*v1.Pipeline, error)
	Update(ctx context.Context, id uint, ns string, body *v1.Pipeline) (*v1.Pipeline, error)
	Delete(ctx context.Context, id uint, ns string, name string) error
}

type pipelineManager struct {
	pipeline repository.PipelineRepository
}

func NewPipelineManager(pipeline repository.PipelineRepository) PipelineManager {
	return &pipelineManager{pipeline: pipeline}
}

// List 列表
func (s *pipelineManager) List(ctx context.Context, id uint, ns string) ([]dto.Pipeline, error) {
	items, err := s.pipeline.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToPipelineDTOs(items), nil
}

// Get 查询
func (s *pipelineManager) Get(ctx context.Context, id uint, ns, name string) (*v1.Pipeline, error) {
	return s.pipeline.Get(ctx, id, ns, name)
}

// Create 创建
func (s *pipelineManager) Create(ctx context.Context, id uint, ns string, body *v1.Pipeline) (*v1.Pipeline, error) {
	return s.pipeline.Create(ctx, id, ns, body)
}

// Update 更新
func (s *pipelineManager) Update(ctx context.Context, id uint, ns string, body *v1.Pipeline) (*v1.Pipeline, error) {
	return s.pipeline.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *pipelineManager) Delete(ctx context.Context, id uint, ns, name string) (err error) {
	return s.pipeline.Delete(ctx, id, ns, name)
}
