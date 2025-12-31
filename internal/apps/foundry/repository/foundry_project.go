package repository

import (
	"context"
	"valyria-backend/internal/apps/foundry/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// ProjectRepository 定义了数据访问层的接口
type ProjectRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Project, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Project, error)
	Create(ctx context.Context, data *model.Project) error
	Update(ctx context.Context, id int, data *model.Project) error
	Delete(ctx context.Context, id int) error
	WithTx(tx *gorm.DB) ProjectRepository
}

// projectRepository 实现了 ProjectRepository 接口
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository 创建新的 ProjectRepository 实例
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

// List 查询列表
func (r *projectRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Project, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *projectRepository) Get(ctx context.Context, id int) (model.Project, error) {
	var data model.Project
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *projectRepository) Create(ctx context.Context, data *model.Project) error {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新
func (r *projectRepository) Update(ctx context.Context, id int, data *model.Project) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除
func (r *projectRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Project{}, id).Error
}

// WithTx 返回一个绑定事务的 Repository
func (r *projectRepository) WithTx(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}
