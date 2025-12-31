package service

import (
	"context"
	"valyria-backend/internal/apps/foundry/model"
	"valyria-backend/internal/apps/foundry/repository"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type CredGitlabManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.CredGitlab, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.CredGitlab, error)
	Create(ctx context.Context, data *model.CredGitlab) error
	Update(ctx context.Context, id int, data *model.CredGitlab) error
	Delete(ctx context.Context, id int) error
}

// credGitlabManager 实现 CredGitlabManager 接口
type credGitlabManager struct {
	repo repository.CredGitlabRepository
	db   *gorm.DB
}

// NewCredGitlabManager 创建新的 CredGitlabManager 实例
func NewCredGitlabManager(repo repository.CredGitlabRepository, db *gorm.DB) CredGitlabManager {
	// 使用 NewGitlabCredentialRepository 函数创建具体实现
	return &credGitlabManager{repo: repo, db: db}
}

// List 列表
func (s *credGitlabManager) List(ctx context.Context, params pg.QueryParams) ([]model.CredGitlab, pg.Pagination, error) {
	return s.repo.List(ctx, params)
}

// Get 查询
func (s *credGitlabManager) Get(ctx context.Context, id int) (model.CredGitlab, error) {
	return s.repo.Get(ctx, id)
}

// Create 创建
func (s *credGitlabManager) Create(ctx context.Context, data *model.CredGitlab) error {
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *credGitlabManager) Update(ctx context.Context, id int, data *model.CredGitlab) error {
	return s.repo.Update(ctx, id, data)
}

// Delete 更新
func (s *credGitlabManager) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
