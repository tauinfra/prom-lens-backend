package repository

import (
	"context"
	"errors"
	"valyria-backend/internal/apps/auth/model"
	pg "valyria-backend/internal/core/pagination"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepository 定义接口
type UserRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.User, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.User, error)
	Create(ctx context.Context, data *model.User) error
	Update(ctx context.Context, id int, data *model.User) (err error)
	Delete(ctx context.Context, id int) error
	Login(ctx context.Context, username, password string) (data model.User, err error)
	UpdatePassword(ctx context.Context, username, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, username, newPassword string) error
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
func (r *userRepository) Get(ctx context.Context, id int) (model.User, error) {
	var data model.User
	if err := r.db.WithContext(ctx).Preload("Modules").First(&data, id).Error; err != nil {
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
	return r.db.WithContext(ctx).Model(&model.User{}).Create(data).Error
}

// Update 更新
func (r *userRepository) Update(ctx context.Context, id int, data *model.User) (err error) {
	var (
		user    model.User
		modules []model.Module
		roles   []model.Role
	)
	// 0. 查出 user（必须，否则 Association 会报错）
	if err = r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return err
	}
	// 1. 查询需要绑定的模块
	if len(data.ModuleIDs) > 0 {
		if err = r.db.WithContext(ctx).Where("id IN ?", data.ModuleIDs).Find(&modules).Error; err != nil {
			return err
		}
	}
	if len(data.RoleIDs) > 0 {
		if err = r.db.WithContext(ctx).Where("id IN ?", data.RoleIDs).Find(&roles).Error; err != nil {
			return err
		}
	}
	// 2. 更新用户(排除密码字段)
	if err = r.db.WithContext(ctx).Model(&user).Omit("Password").Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	// 3. 替换用户模块（核心逻辑）
	if err = r.db.WithContext(ctx).Model(&user).Association("Modules").Replace(modules); err != nil {
		return err
	}

	// 4. 替换用户模块（核心逻辑）
	if err = r.db.WithContext(ctx).Model(&user).Association("Roles").Replace(roles); err != nil {
		return err
	}
	// 查询更新后的最新用户信息
	return r.db.WithContext(ctx).First(&data, id).Error
}

// Delete 删除
func (r *userRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// Login 登录
func (r *userRepository) Login(ctx context.Context, username, password string) (data model.User, err error) {
	// 统一错误提示，防止用户枚举
	const authErrorMsg = "账号或密码错误，请重新输入"
	if err = r.db.WithContext(ctx).Preload("Roles").Where("username = ?", username).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New(authErrorMsg)
		}
	}
	// 非活跃用户禁止登录
	if data.IsActive != nil && !*data.IsActive {
		return data, errors.New("该账号已禁用，请联系管理员！")
	}
	// 密码校验
	if err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(password)); err != nil {
		return data, errors.New(authErrorMsg)
	}
	return
}

// UpdatePassword 修改密码
func (r *userRepository) UpdatePassword(ctx context.Context, username, oldPassword, newPassword string) error {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}
	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return err
	}
	// 修改新秘密
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Update("password", string(hashPassword)).Error; err != nil {
		return err
	}
	return nil
}

// ResetPassword 重置密码
func (r *userRepository) ResetPassword(ctx context.Context, username, newPassword string) error {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}
	// 修改新秘密
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Update("password", string(hashPassword)).Error; err != nil {
		return err
	}
	return nil
}

// WithTx 返回一个绑定事务的 Repository
func (r *userRepository) WithTx(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
