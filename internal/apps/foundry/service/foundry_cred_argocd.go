package service

import (
	"context"
	"valyria-backend/internal/apps/foundry/model"
	"valyria-backend/internal/apps/foundry/repository"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type CredArgoCDManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.CredArgoCD, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.CredArgoCD, error)
	Create(ctx context.Context, data *model.CredArgoCD) error
	Update(ctx context.Context, id int, data *model.CredArgoCD) error
	Delete(ctx context.Context, id int) error
}

// CredArgoCDManager 实现 CredArgoCDManager 接口
type credArgoCDManager struct {
	repo repository.CredArgoCDRepository
	db   *gorm.DB
}

// NewCredArgoCDManager 创建新的 CredArgoCDManager 实例
func NewCredArgoCDManager(repo repository.CredArgoCDRepository, db *gorm.DB) CredArgoCDManager {
	// 使用 NewGitlabCredentialRepository 函数创建具体实现
	return &credArgoCDManager{repo: repo, db: db}
}

// List 列表
func (s *credArgoCDManager) List(ctx context.Context, params pg.QueryParams) ([]model.CredArgoCD, pg.Pagination, error) {
	return s.repo.List(ctx, params)
}

// Get 查询
func (s *credArgoCDManager) Get(ctx context.Context, id int) (model.CredArgoCD, error) {
	return s.repo.Get(ctx, id)
}

// Create 创建
func (s *credArgoCDManager) Create(ctx context.Context, data *model.CredArgoCD) error {
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *credArgoCDManager) Update(ctx context.Context, id int, data *model.CredArgoCD) error {
	return s.repo.Update(ctx, id, data)
}

// Delete 更新
func (s *credArgoCDManager) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
