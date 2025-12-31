package service

import (
	"context"
	"valyria-backend/internal/apps/foundry/model"
	"valyria-backend/internal/apps/foundry/repository"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type PipelineManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Pipeline, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Pipeline, error)
	Create(ctx context.Context, data *model.Pipeline) error
	Update(ctx context.Context, id int, data *model.Pipeline) error
	Delete(ctx context.Context, id int) error
}

// pipelineManager 实现 PipelineManager 接口
type pipelineManager struct {
	repo repository.PipelineRepository
	db   *gorm.DB
}

// NewPipelineManager 创建新的 PipelineManager 实例
func NewPipelineManager(repo repository.PipelineRepository, db *gorm.DB) PipelineManager {
	// 使用 NewGitlabCredentialRepository 函数创建具体实现
	return &pipelineManager{repo: repo, db: db}
}

// List 列表
func (s *pipelineManager) List(ctx context.Context, params pg.QueryParams) ([]model.Pipeline, pg.Pagination, error) {
	return s.repo.List(ctx, params)
}

// Get 查询
func (s *pipelineManager) Get(ctx context.Context, id int) (model.Pipeline, error) {
	return s.repo.Get(ctx, id)
}

// Create 创建
func (s *pipelineManager) Create(ctx context.Context, data *model.Pipeline) error {
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *pipelineManager) Update(ctx context.Context, id int, data *model.Pipeline) error {
	return s.repo.Update(ctx, id, data)
}

// Delete 更新
func (s *pipelineManager) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
