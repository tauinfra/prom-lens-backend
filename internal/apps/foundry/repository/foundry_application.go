package repository

import (
	"context"
	"valyria-backend/internal/apps/foundry/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// ApplicationRepository 定义了数据访问层的接口
type ApplicationRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Application, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Application, error)
	Create(ctx context.Context, data *model.Application) error
	Update(ctx context.Context, id int, data *model.Application) error
	Delete(ctx context.Context, id int) error
	WithTx(tx *gorm.DB) ApplicationRepository
}

// applicationRepository 实现了 ApplicationRepository 接口
type applicationRepository struct {
	db *gorm.DB
}

// NewApplicationRepository 创建新的 ApplicationRepository 实例
func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

// List 查询列表
func (r *applicationRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Application, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *applicationRepository) Get(ctx context.Context, id int) (model.Application, error) {
	var data model.Application
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *applicationRepository) Create(ctx context.Context, data *model.Application) error {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新
func (r *applicationRepository) Update(ctx context.Context, id int, data *model.Application) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除
func (r *applicationRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Application{}, id).Error
}

// WithTx 返回一个绑定事务的 Repository
func (r *applicationRepository) WithTx(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}
