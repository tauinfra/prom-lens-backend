package repository

import (
	"context"
	"prom-lens-backend/internal/apps/prometheus/model"
	pg "prom-lens-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// TargetGroupRepository 定义接口
type TargetGroupRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.TargetGroup, pg.Pagination, error)
	HasTargets(ctx context.Context, groupID int) (bool, error)
	Get(ctx context.Context, id int) (model.TargetGroup, error)
	Create(ctx context.Context, data *model.TargetGroup) error
	Update(ctx context.Context, id int, data *model.TargetGroup) error
	Delete(ctx context.Context, id int) error
	WithTx(tx *gorm.DB) TargetGroupRepository
}

// targetGroupRepository 实现了 TargetGroupRepository 接口
type targetGroupRepository struct {
	db *gorm.DB
}

// NewTargetGroupRepository 创建新的 TargetGroupRepository 实例
func NewTargetGroupRepository(db *gorm.DB) TargetGroupRepository {
	return &targetGroupRepository{db: db}
}

// List 列表
func (r *targetGroupRepository) List(ctx context.Context, params pg.QueryParams) (data []model.TargetGroup, pagination pg.Pagination, err error) {
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

func (r *targetGroupRepository) HasTargets(ctx context.Context, groupID int) (bool, error) {
	var exists bool
	if err := r.db.WithContext(ctx).
		Raw("SELECT EXISTS(SELECT 1 FROM prom_lens_prometheus_target WHERE group_id = ?)", groupID).
		Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

// Get 查询
func (r *targetGroupRepository) Get(ctx context.Context, id int) (data model.TargetGroup, err error) {
	if err = r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return
}

// Create 创建
func (r *targetGroupRepository) Create(ctx context.Context, data *model.TargetGroup) error {
	return r.db.WithContext(ctx).Create(data).Error
}

// Update 更新
func (r *targetGroupRepository) Update(ctx context.Context, id int, data *model.TargetGroup) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error
}

// Delete 删除
func (r *targetGroupRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.TargetGroup{}).Error
}

func (r *targetGroupRepository) WithTx(db *gorm.DB) TargetGroupRepository {
	return &targetGroupRepository{db: db}
}
