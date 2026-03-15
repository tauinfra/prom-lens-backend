package repository

import (
	"context"
	"valyria-backend/internal/apps/authn/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// MenuRepository 定义接口
type MenuRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Menu, pg.Pagination, error)
	ListAll(ctx context.Context) ([]model.Menu, error)
	Get(ctx context.Context, menuID uint) (model.Menu, error)
	Create(ctx context.Context, data *model.Menu) error
	Update(ctx context.Context, menuID uint, data *model.Menu) (err error)
	Delete(ctx context.Context, menuID uint) error
	WithTx(tx *gorm.DB) MenuRepository
}

// menuRepository 实现了 MenuRepository 接口
type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

// List 列表
func (r *menuRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Menu, pagination pg.Pagination, err error) {
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// ListAll 全量列表（用于构建树，按 rank 排序）
func (r *menuRepository) ListAll(ctx context.Context) ([]model.Menu, error) {
	var data []model.Menu
	err := r.db.WithContext(ctx).Order("`rank` ASC, id ASC").Find(&data).Error
	return data, err
}

// Get 查询
func (r *menuRepository) Get(ctx context.Context, id uint) (model.Menu, error) {
	var data model.Menu
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *menuRepository) Create(ctx context.Context, data *model.Menu) error {
	return r.db.WithContext(ctx).Create(data).Error
}

// Update 更新
func (r *menuRepository) Update(ctx context.Context, id uint, data *model.Menu) (err error) {
	if err = r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).First(&data, id).Error
}

// Delete 删除
func (r *menuRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Menu{}, id).Error
}

// WithTx 返回一个绑定事务的 Repository
func (r *menuRepository) WithTx(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}
