package service

import (
	"context"
	"valyria-backend/internal/apps/foundry/model"
	"valyria-backend/internal/apps/foundry/repository"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type EnvironmentManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Environment, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Environment, error)
	Create(ctx context.Context, data *model.Environment) error
	Update(ctx context.Context, id int, data *model.Environment) error
	Delete(ctx context.Context, id int) error
}

// EnvironmentManager 实现 EnvironmentManager 接口
type environmentManager struct {
	repo repository.EnvironmentRepository
	db   *gorm.DB
}

// NewEnvironmentManager 创建新的 EnvironmentManager 实例
func NewEnvironmentManager(repo repository.EnvironmentRepository, db *gorm.DB) EnvironmentManager {
	// 使用 NewGitlabCredentialRepository 函数创建具体实现
	return &environmentManager{repo: repo, db: db}
}

// List 列表
func (s *environmentManager) List(ctx context.Context, params pg.QueryParams) ([]model.Environment, pg.Pagination, error) {
	return s.repo.List(ctx, params)
}

// Get 查询
func (s *environmentManager) Get(ctx context.Context, id int) (model.Environment, error) {
	return s.repo.Get(ctx, id)
}

// Create 创建
func (s *environmentManager) Create(ctx context.Context, data *model.Environment) error {
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *environmentManager) Update(ctx context.Context, id int, data *model.Environment) error {
	return s.repo.Update(ctx, id, data)
}

// Delete 更新
func (s *environmentManager) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
