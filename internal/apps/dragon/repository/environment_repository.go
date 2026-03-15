package repository

import (
	"context"
	"valyria-backend/internal/apps/dragon/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// EnvironmentRepository 定义了数据访问层的接口
type EnvironmentRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Environment, pg.Pagination, error)
	Get(ctx context.Context, projectID, id uint) (model.Environment, error)
	Create(ctx context.Context, data *model.Environment) error
	Update(ctx context.Context, projectID, id uint, data *model.Environment) error
	Delete(ctx context.Context, projectID, id uint) error
	WithTx(tx *gorm.DB) EnvironmentRepository
}

// environmentRepository 实现了 EnvironmentRepository 接口
type environmentRepository struct {
	db *gorm.DB
}

// NewEnvironmentRepository 创建新的 EnvironmentRepository 实例
func NewEnvironmentRepository(db *gorm.DB) EnvironmentRepository {
	return &environmentRepository{db: db}
}

// List 查询列表
func (r *environmentRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Environment, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *environmentRepository) Get(ctx context.Context, projectID, id uint) (model.Environment, error) {
	var data model.Environment
	if err := r.db.WithContext(ctx).
		Preload("Project").
		Preload("Harbor").
		Where("id = ? AND project_id = ?", id, projectID).
		First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *environmentRepository) Create(ctx context.Context, data *model.Environment) error {
	return r.db.WithContext(ctx).Create(data).Error
}

// Update 更新
func (r *environmentRepository) Update(ctx context.Context, projectID, id uint, data *model.Environment) error {
	return r.db.WithContext(ctx).Where("id = ? AND project_id = ?", id, projectID).Updates(data).Error
}

// Delete 删除
func (r *environmentRepository) Delete(ctx context.Context, projectID, id uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND project_id = ?", id, projectID).
		Delete(&model.Environment{}, id).Error
}

// WithTx 返回一个绑定事务的 Repository
func (r *environmentRepository) WithTx(db *gorm.DB) EnvironmentRepository {
	return &environmentRepository{db: db}
}
