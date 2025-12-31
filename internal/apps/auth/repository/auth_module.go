package repository

import (
	"context"
	"valyria-backend/internal/apps/auth/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// ModuleRepository 定义接口
type ModuleRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Module, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Module, error)
	Create(ctx context.Context, data *model.Module) error
	Update(ctx context.Context, id int, data *model.Module) (err error)
	Delete(ctx context.Context, id int) error
	WithTx(tx *gorm.DB) ModuleRepository // 提供一个事务操作接口
}

// moduleRepository 实现了 ModuleRepository 接口
type moduleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) ModuleRepository {
	return &moduleRepository{db: db}
}

// List 列表
func (r *moduleRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Module, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *moduleRepository) Get(ctx context.Context, id int) (model.Module, error) {
	var data model.Module
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *moduleRepository) Create(ctx context.Context, data *model.Module) error {
	// 创建用户
	return r.db.WithContext(ctx).Create(data).Error
}

// Update 更新
func (r *moduleRepository) Update(ctx context.Context, id int, data *model.Module) (err error) {
	if err = r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	// 查询更新后的最新用户信息
	return r.db.WithContext(ctx).First(&data, id).Error
}

// Delete 删除
func (r *moduleRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Module{}, id).Error
}

// WithTx 返回一个绑定事务的 Repository
func (r *moduleRepository) WithTx(db *gorm.DB) ModuleRepository {
	return &moduleRepository{db: db}
}
