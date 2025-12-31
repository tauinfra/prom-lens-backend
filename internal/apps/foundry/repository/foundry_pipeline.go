package repository

import (
	"context"
	"valyria-backend/internal/apps/foundry/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// PipelineRepository 定义了数据访问层的接口
type PipelineRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Pipeline, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Pipeline, error)
	Create(ctx context.Context, data *model.Pipeline) error
	Update(ctx context.Context, id int, data *model.Pipeline) error
	Delete(ctx context.Context, id int) error
	WithTx(tx *gorm.DB) PipelineRepository
}

// pipelineRepository 实现了 PipelineRepository 接口
type pipelineRepository struct {
	db *gorm.DB
}

// NewPipelineRepository 创建新的 PipelineRepository 实例
func NewPipelineRepository(db *gorm.DB) PipelineRepository {
	return &pipelineRepository{db: db}
}

// List 查询列表
func (r *pipelineRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Pipeline, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *pipelineRepository) Get(ctx context.Context, id int) (model.Pipeline, error) {
	var data model.Pipeline
	if err := r.db.WithContext(ctx).
		Preload("Environment.Project").
		Preload("Environment.CredHarbor").
		Preload("Application").
		First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *pipelineRepository) Create(ctx context.Context, data *model.Pipeline) error {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新
func (r *pipelineRepository) Update(ctx context.Context, id int, data *model.Pipeline) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除
func (r *pipelineRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Pipeline{}, id).Error
}

// WithTx 返回一个绑定事务的 Repository
func (r *pipelineRepository) WithTx(db *gorm.DB) PipelineRepository {
	return &pipelineRepository{db: db}
}
