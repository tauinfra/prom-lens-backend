package service

import (
	"context"
	"valyria-backend/internal/apps/auth/model"
	"valyria-backend/internal/apps/auth/repository"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type RoleManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Role, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Role, error)
	Create(ctx context.Context, data *model.Role) error
	Update(ctx context.Context, id int, data *model.Role) error
	Delete(ctx context.Context, id int) error
}

// roleManager 实现 RoleManager 接口
type roleManager struct {
	repo repository.RoleRepository
	db   *gorm.DB
}

// NewRoleManager 创建新的 RoleManager 实例
func NewRoleManager(db *gorm.DB, repo repository.RoleRepository) RoleManager {
	return &roleManager{db: db, repo: repo}
}

// List 列表
func (s *roleManager) List(ctx context.Context, params pg.QueryParams) ([]model.Role, pg.Pagination, error) {
	return s.repo.List(ctx, params)
}

// Get 查询
func (s *roleManager) Get(ctx context.Context, id int) (model.Role, error) {
	return s.repo.Get(ctx, id)
}

// Create 创建
func (s *roleManager) Create(ctx context.Context, data *model.Role) error {
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *roleManager) Update(ctx context.Context, id int, data *model.Role) error {
	return s.repo.Update(ctx, id, data)
}

// Delete 删除
func (s *roleManager) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
