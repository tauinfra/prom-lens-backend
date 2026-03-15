package repository

import (
	"context"
	"valyria-backend/internal/apps/dragon/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// PipelineRepository 定义了数据访问层的接口
type PipelineRepository interface {
	List(ctx context.Context, projectID, environmentID uint, params pg.QueryParams) ([]model.Pipeline, pg.Pagination, error)
	Get(ctx context.Context, projectID, environmentID, pipelineID uint) (model.Pipeline, error)
	Create(ctx context.Context, data *model.Pipeline) error
	Update(ctx context.Context, projectID, environmentID, pipelineID uint, data *model.Pipeline) error
	Delete(ctx context.Context, projectID, environmentID, pipelineID uint) error
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
func (r *pipelineRepository) List(ctx context.Context, projectID, environmentID uint, params pg.QueryParams) (data []model.Pipeline, pagination pg.Pagination, err error) {
	query := r.db.WithContext(ctx).
		Model(&model.Pipeline{}).
		Joins("JOIN valyria_dragon_environment e ON e.id = valyria_dragon_pipeline.environment_id").
		Where("valyria_dragon_pipeline.environment_id = ? AND e.project_id = ?", environmentID, projectID)
	// 分页查询
	if pagination, err = pg.Paginate(query, &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *pipelineRepository) Get(ctx context.Context, projectID, environmentID, pipelineID uint) (model.Pipeline, error) {
	var data model.Pipeline
	if err := r.db.WithContext(ctx).
		Preload("Environment.Project").
		Preload("Environment.Harbor").
		Preload("Gitlab").
		Joins("JOIN valyria_dragon_environment e ON e.id = valyria_dragon_pipeline.environment_id").
		Where("valyria_dragon_pipeline.id = ? AND valyria_dragon_pipeline.environment_id = ? AND e.project_id = ?", pipelineID, environmentID, projectID).
		First(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *pipelineRepository) Create(ctx context.Context, data *model.Pipeline) error {
	return r.db.WithContext(ctx).Create(data).Error
}

// Update 更新
func (r *pipelineRepository) Update(ctx context.Context, projectID, environmentID, pipelineID uint, data *model.Pipeline) error {
	_, err := r.Get(ctx, projectID, environmentID, pipelineID)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ? AND environment_id = ?", pipelineID, environmentID).Updates(data).Error
}

// Delete 删除
func (r *pipelineRepository) Delete(ctx context.Context, projectID, environmentID, pipelineID uint) error {
	if _, err := r.Get(ctx, projectID, environmentID, pipelineID); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(&model.Pipeline{}, pipelineID).Error
}

// WithTx 返回一个绑定事务的 Repository
func (r *pipelineRepository) WithTx(db *gorm.DB) PipelineRepository {
	return &pipelineRepository{db: db}
}
