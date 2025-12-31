package repository

import (
	"context"
	"valyria-backend/internal/apps/auth/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// RoleRepository 定义接口
type RoleRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Role, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Role, error)
	Create(ctx context.Context, data *model.Role) error
	Update(ctx context.Context, id int, data *model.Role) (err error)
	Delete(ctx context.Context, id int) error
	WithTx(tx *gorm.DB) ModuleRepository // 提供一个事务操作接口
}

// roleRepository 实现了 RoleRepository 接口
type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// List 列表
func (r *roleRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Role, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *roleRepository) Get(ctx context.Context, id int) (model.Role, error) {
	var data model.Role
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *roleRepository) Create(ctx context.Context, data *model.Role) error {
	// 创建用户
	return r.db.WithContext(ctx).Create(data).Error
}

// Update 更新
func (r *roleRepository) Update(ctx context.Context, id int, data *model.Role) (err error) {
	if err = r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	// 查询更新后的最新用户信息
	return r.db.WithContext(ctx).First(&data, id).Error
}

// Delete 删除
func (r *roleRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Role{}, id).Error
}

// WithTx 返回一个绑定事务的 Repository
func (r *roleRepository) WithTx(db *gorm.DB) ModuleRepository {
	return &moduleRepository{db: db}
}
