package service

import (
	"context"
	"valyria-backend/internal/apps/foundry/model"
	"valyria-backend/internal/apps/foundry/repository"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type ApplicationManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Application, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Application, error)
	Create(ctx context.Context, data *model.Application) error
	Update(ctx context.Context, id int, data *model.Application) error
	Delete(ctx context.Context, id int) error
}

// applicationManager 实现 ApplicationManager 接口
type applicationManager struct {
	repo repository.ApplicationRepository
	db   *gorm.DB
}

// NewApplicationManager 创建新的 ApplicationManager 实例
func NewApplicationManager(repo repository.ApplicationRepository, db *gorm.DB) ApplicationManager {
	// 使用 NewGitlabCredentialRepository 函数创建具体实现
	return &applicationManager{repo: repo, db: db}
}

// List 列表
func (s *applicationManager) List(ctx context.Context, params pg.QueryParams) ([]model.Application, pg.Pagination, error) {
	return s.repo.List(ctx, params)
}

// Get 查询
func (s *applicationManager) Get(ctx context.Context, id int) (model.Application, error) {
	return s.repo.Get(ctx, id)
}

// Create 创建
func (s *applicationManager) Create(ctx context.Context, data *model.Application) error {
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *applicationManager) Update(ctx context.Context, id int, data *model.Application) error {
	return s.repo.Update(ctx, id, data)
}

// Delete 更新
func (s *applicationManager) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
