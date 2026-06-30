package repository

import (
	"context"
	"prom-lens-backend/internal/apps/prometheus/model"
	pg "prom-lens-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// RecordRepository 定义接口
type RecordRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Record, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Record, error)
	GetByGroupAndName(ctx context.Context, groupID int, name string) (model.Record, error)
	Create(ctx context.Context, data *model.Record) error
	Update(ctx context.Context, id int, data *model.Record) error
	Delete(ctx context.Context, id int) error
	WithTx(tx *gorm.DB) RecordRepository
}

// recordRepository 实现了 RecordRepository 接口
type recordRepository struct {
	db *gorm.DB
}

// NewRecordRepository 创建新的 RecordRepository 实例
func NewRecordRepository(db *gorm.DB) RecordRepository {
	return &recordRepository{db: db}
}

// List 列表
func (r *recordRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Record, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *recordRepository) Get(ctx context.Context, id int) (data model.Record, err error) {
	if err = r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return
}

// GetByGroupAndName 按组与记录名查询
func (r *recordRepository) GetByGroupAndName(ctx context.Context, groupID int, name string) (data model.Record, err error) {
	err = r.db.WithContext(ctx).Where("group_id = ? AND name = ?", groupID, name).First(&data).Error
	return
}

// Create 创建
func (r *recordRepository) Create(ctx context.Context, data *model.Record) error {
	return r.db.WithContext(ctx).Create(data).Error
}

// Update 更新
func (r *recordRepository) Update(ctx context.Context, id int, data *model.Record) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error
}

// Delete 删除
func (r *recordRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Record{}).Error
}

func (r *recordRepository) WithTx(db *gorm.DB) RecordRepository {
	return &recordRepository{db: db}
}
