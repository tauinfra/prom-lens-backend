package repository

import (
	"context"
	"valyria-backend/internal/apps/prometheus/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// TargetRepository 定义接口
type TargetRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Target, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Target, error)
	Create(ctx context.Context, data *model.Target) error
	Update(ctx context.Context, id int, data *model.Target) error
	Delete(ctx context.Context, id int) error
	WithTx(tx *gorm.DB) TargetRepository
}

// targetRepository 实现了 TargetRepository 接口
type targetRepository struct {
	db *gorm.DB
}

// NewTargetRepository 创建新的 TargetRepository 实例
func NewTargetRepository(db *gorm.DB) TargetRepository {
	return &targetRepository{db: db}
}

// List 列表
func (r *targetRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Target, pagination pg.Pagination, err error) {
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *targetRepository) Get(ctx context.Context, id int) (data model.Target, err error) {
	if err = r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return
}

// Create 创建
func (r *targetRepository) Create(ctx context.Context, data *model.Target) error {
	return r.db.WithContext(ctx).Create(data).Error
}

// Update 更新
func (r *targetRepository) Update(ctx context.Context, id int, data *model.Target) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error
}

// Delete 删除
func (r *targetRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Target{}).Error
}

func (r *targetRepository) WithTx(db *gorm.DB) TargetRepository {
	return &targetRepository{db: db}
}
