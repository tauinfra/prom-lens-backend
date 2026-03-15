package repository

import (
	"context"
	"valyria-backend/internal/apps/authn/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// RoleRepository 定义接口
type RoleRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Role, pg.Pagination, error)
	Get(ctx context.Context, roleID uint) (model.Role, error)
	Create(ctx context.Context, data *model.Role) error
	Update(ctx context.Context, roleID uint, data *model.Role) (err error)
	Delete(ctx context.Context, roleID uint) error
	GetRolePermissions(ctx context.Context, roleID uint) ([]model.Permission, error)
	UpdateRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
	WithTx(tx *gorm.DB) RoleRepository // 提供一个事务操作接口
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
func (r *roleRepository) Get(ctx context.Context, id uint) (model.Role, error) {
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
func (r *roleRepository) Update(ctx context.Context, id uint, data *model.Role) (err error) {
	if err = r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	// 查询更新后的最新用户信息
	return r.db.WithContext(ctx).First(&data, id).Error
}

// Delete 删除
func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Role{}, id).Error
}

// GetRolePermissions 查询角色权限
func (r *roleRepository) GetRolePermissions(ctx context.Context, roleID uint) (data []model.Permission, err error) {
	err = r.db.WithContext(ctx).Table("valyria_authn_permission p").
		Select("p.id, p.code, p.name, p.module").
		Joins("join valyria_authn_role_permission rp on rp.permission_id = p.id").
		Where("rp.role_id = ?", roleID).
		Scan(&data).Error
	return data, err
}

// UpdateRolePermissions 更新角色权限
func (r *roleRepository) UpdateRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	var data []model.RolePermission
	tx := r.db.WithContext(ctx).Begin()
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()
	if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
		return err
	}
	for _, permissionID := range permissionIDs {
		data = append(data, model.RolePermission{RoleID: roleID, PermissionID: permissionID})
	}
	if len(data) > 0 {
		if err := tx.Create(&data).Error; err != nil {
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	committed = true
	return nil
}

// WithTx 返回一个绑定事务的 Repository
func (r *roleRepository) WithTx(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}
