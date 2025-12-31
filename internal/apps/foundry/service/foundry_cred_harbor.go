package service

import (
	"context"
	"valyria-backend/internal/apps/foundry/model"
	"valyria-backend/internal/apps/foundry/repository"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type CredHarborManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.CredHarbor, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.CredHarbor, error)
	Create(ctx context.Context, data *model.CredHarbor) error
	Update(ctx context.Context, id int, data *model.CredHarbor) error
	Delete(ctx context.Context, id int) error
}

// CredHarborManager 实现 CredHarborManager 接口
type credHarborManager struct {
	repo repository.CredHarborRepository
	db   *gorm.DB
}

// NewCredHarborManager 创建新的 CredHarborManager 实例
func NewCredHarborManager(repo repository.CredHarborRepository, db *gorm.DB) CredHarborManager {
	// 使用 NewGitlabCredentialRepository 函数创建具体实现
	return &credHarborManager{repo: repo, db: db}
}

// List 列表
func (s *credHarborManager) List(ctx context.Context, params pg.QueryParams) ([]model.CredHarbor, pg.Pagination, error) {
	return s.repo.List(ctx, params)
}

// Get 查询
func (s *credHarborManager) Get(ctx context.Context, id int) (model.CredHarbor, error) {
	return s.repo.Get(ctx, id)
}

// Create 创建
func (s *credHarborManager) Create(ctx context.Context, data *model.CredHarbor) error {
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *credHarborManager) Update(ctx context.Context, id int, data *model.CredHarbor) error {
	return s.repo.Update(ctx, id, data)
}

// Delete 更新
func (s *credHarborManager) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
