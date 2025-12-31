package service

import (
	"context"
	"valyria-backend/internal/apps/auth/model"
	"valyria-backend/internal/apps/auth/repository"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type ModuleManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Module, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Module, error)
	Create(ctx context.Context, data *model.Module) error
	Update(ctx context.Context, id int, data *model.Module) error
	Delete(ctx context.Context, id int) error
}

// moduleManager 实现 ModuleManager 接口
type moduleManager struct {
	repo repository.ModuleRepository
	db   *gorm.DB
}

// NewModuleManager 创建新的 ModuleManager 实例
func NewModuleManager(db *gorm.DB, repo repository.ModuleRepository) ModuleManager {
	return &moduleManager{db: db, repo: repo}
}

// List 列表
func (s *moduleManager) List(ctx context.Context, params pg.QueryParams) ([]model.Module, pg.Pagination, error) {
	return s.repo.List(ctx, params)
}

// Get 查询
func (s *moduleManager) Get(ctx context.Context, id int) (model.Module, error) {
	return s.repo.Get(ctx, id)
}

// Create 创建
func (s *moduleManager) Create(ctx context.Context, data *model.Module) error {
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *moduleManager) Update(ctx context.Context, id int, data *model.Module) error {
	return s.repo.Update(ctx, id, data)
}

// Delete 删除
func (s *moduleManager) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
