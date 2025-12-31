package service

import (
	"context"
	"valyria-backend/internal/apps/foundry/model"
	"valyria-backend/internal/apps/foundry/repository"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type ProjectManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Project, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Project, error)
	Create(ctx context.Context, data *model.Project) error
	Update(ctx context.Context, id int, data *model.Project) error
	Delete(ctx context.Context, id int) error
}

// projectManager 实现 ProjectManager 接口
type projectManager struct {
	repo repository.ProjectRepository
	db   *gorm.DB
}

// NewProjectManager 创建新的 ProjectManager 实例
func NewProjectManager(repo repository.ProjectRepository, db *gorm.DB) ProjectManager {
	// 使用 NewGitlabCredentialRepository 函数创建具体实现
	return &projectManager{repo: repo, db: db}
}

// List 列表
func (s *projectManager) List(ctx context.Context, params pg.QueryParams) (data []model.Project, pg pg.Pagination, err error) {

	return s.repo.List(ctx, params)
}

// Get 查询
func (s *projectManager) Get(ctx context.Context, id int) (model.Project, error) {
	return s.repo.Get(ctx, id)
}

// Create 创建
func (s *projectManager) Create(ctx context.Context, data *model.Project) error {
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *projectManager) Update(ctx context.Context, id int, data *model.Project) error {
	return s.repo.Update(ctx, id, data)
}

// Delete 更新
func (s *projectManager) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
