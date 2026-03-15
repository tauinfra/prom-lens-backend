package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/apps/authn/model"
	pg "valyria-backend/internal/core/pagination"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepository 定义接口
type UserRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.User, pg.Pagination, error)
	Get(ctx context.Context, id uint) (model.User, error)
	Create(ctx context.Context, data *model.User) error
	Update(ctx context.Context, id uint, data *model.User) (err error)
	Delete(ctx context.Context, id uint) error
	ResetPassword(ctx context.Context, id uint, password string) error
	GetUserRoles(ctx context.Context, id uint) ([]model.Role, error)
	UpdateUserRoles(ctx context.Context, id uint, roleIDs []uint) error
	GetUserMenus(ctx context.Context, id uint) ([]model.Menu, error)
	UpdateUserMenus(ctx context.Context, id uint, menuIDs []uint) error
	WithTx(tx *gorm.DB) UserRepository // 提供一个事务操作接口
}

// userRepository 实现了 UserRepository 接口
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// List 列表
func (r *userRepository) List(ctx context.Context, params pg.QueryParams) (data []model.User, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *userRepository) Get(ctx context.Context, id uint) (data model.User, err error) {
	if err = r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *userRepository) Create(ctx context.Context, data *model.User) error {
	// 密码字段加密
	hash, _ := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	data.Password = string(hash)
	// 创建用户
	return r.db.WithContext(ctx).Create(data).Error
}

// Update 更新
func (r *userRepository) Update(ctx context.Context, id uint, data *model.User) (err error) {
	var (
		user model.User
	)
	return r.db.WithContext(ctx).Model(&user).Where("id = ?", id).Updates(data).Error
}

// Delete 删除
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// ResetPassword 重置密码
func (r *userRepository) ResetPassword(ctx context.Context, id uint, password string) error {
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("password", password).Error; err != nil {
		return err
	}
	return nil
}

// GetUserRoles 查询用户的所有角色
func (r *userRepository) GetUserRoles(ctx context.Context, id uint) (data []model.Role, err error) {
	err = r.db.WithContext(ctx).
		Table("valyria_authn_role r").
		Joins("JOIN valyria_authn_user_role ur ON ur.role_id = r.id").
		Where("ur.user_id = ?", id).
		Scan(&data).Error
	return data, err
}

func (r *userRepository) UpdateUserRoles(ctx context.Context, id uint, roleIDs []uint) error {
	var data []model.UserRole
	tx := r.db.WithContext(ctx).Begin()
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()
	if err := tx.Where("user_id = ?", id).Delete(&model.UserRole{}).Error; err != nil {
		return err
	}
	for _, roleID := range roleIDs {
		data = append(data, model.UserRole{UserID: id, RoleID: roleID})
	}
	if len(data) > 0 {
		fmt.Println(data)
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

// GetUserMenus 查询用户的所有菜单
func (r *userRepository) GetUserMenus(ctx context.Context, id uint) (data []model.Menu, err error) {
	err = r.db.WithContext(ctx).
		Table("valyria_authn_menu m").
		Joins("JOIN valyria_authn_user_menu um ON um.menu_id = m.id").
		Where("um.user_id = ?", id).
		Order("m.`rank` ASC, m.id ASC").
		Scan(&data).Error
	return data, err
}

// UpdateUserMenus 更新用户菜单（先删后插）
func (r *userRepository) UpdateUserMenus(ctx context.Context, id uint, menuIDs []uint) error {
	tx := r.db.WithContext(ctx).Begin()
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()
	if err := tx.Where("user_id = ?", id).Delete(&model.UserMenu{}).Error; err != nil {
		return err
	}
	var data []model.UserMenu
	for _, menuID := range menuIDs {
		data = append(data, model.UserMenu{UserID: id, MenuID: menuID})
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
func (r *userRepository) WithTx(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
