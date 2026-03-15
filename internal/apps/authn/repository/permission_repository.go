package repository

import (
	"context"
	"valyria-backend/internal/apps/authn/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// PermissionRepository 定义接口
type PermissionRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Permission, pg.Pagination, error)
	Get(ctx context.Context, id uint) (model.Permission, error)
	Create(ctx context.Context, data *model.Permission) error
	Update(ctx context.Context, id uint, data *model.Permission) (err error)
	Delete(ctx context.Context, id uint) error
	WithTx(tx *gorm.DB) PermissionRepository // 提供一个事务操作接口
}

// permissionRepository 实现了 PermissionRepository 接口
type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

// List 列表
func (r *permissionRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Permission, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *permissionRepository) Get(ctx context.Context, id uint) (model.Permission, error) {
	var data model.Permission
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *permissionRepository) Create(ctx context.Context, data *model.Permission) error {
	// 创建用户
	return r.db.WithContext(ctx).Create(data).Error
}

// Update 更新
func (r *permissionRepository) Update(ctx context.Context, id uint, data *model.Permission) (err error) {
	if err = r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	// 查询更新后的最新用户信息
	return r.db.WithContext(ctx).First(&data, id).Error
}

// Delete 删除
func (r *permissionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Permission{}, id).Error
}

// WithTx 返回一个绑定事务的 Repository
func (r *permissionRepository) WithTx(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}
