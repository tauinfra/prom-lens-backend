package repository

import (
	"context"
	"valyria-backend/internal/apps/prometheus/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// GroupRepository 定义接口
type GroupRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Group, pg.Pagination, error)
	HasRules(ctx context.Context, groupID int) (bool, error)
	HasRecords(ctx context.Context, groupID int) (bool, error)
	Get(ctx context.Context, id int) (model.Group, error)
	Create(ctx context.Context, data *model.Group) error
	Update(ctx context.Context, id int, data *model.Group) error
	Delete(ctx context.Context, id int) error
	WithTx(tx *gorm.DB) GroupRepository
}

// groupRepository 实现了 GroupRepository 接口
type groupRepository struct {
	db *gorm.DB
}

// NewGroupRepository 创建新的 GroupRepository 实例
func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}

// List 列表
func (r *groupRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Group, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

func (r *groupRepository) HasRules(ctx context.Context, groupID int) (bool, error) {
	var exists bool
	if err := r.db.WithContext(ctx).
		Raw("SELECT EXISTS(SELECT 1 FROM valyria_prometheus_rule WHERE group_id = ?)", groupID).
		Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

func (r *groupRepository) HasRecords(ctx context.Context, groupID int) (bool, error) {
	var exists bool
	if err := r.db.WithContext(ctx).
		Raw("SELECT EXISTS(SELECT 1 FROM valyria_prometheus_record WHERE group_id = ?)", groupID).
		Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

// Get 查询
func (r *groupRepository) Get(ctx context.Context, id int) (data model.Group, err error) {
	if err = r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return
}

// Create 创建
func (r *groupRepository) Create(ctx context.Context, data *model.Group) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *groupRepository) Update(ctx context.Context, id int, data *model.Group) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error
}

func (r *groupRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Group{}).Error
}

func (r *groupRepository) WithTx(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}
