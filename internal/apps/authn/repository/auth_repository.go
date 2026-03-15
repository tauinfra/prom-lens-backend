package repository

import (
	"context"
	"errors"
	"time"
	"valyria-backend/internal/apps/authn/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthRepository 定义接口
type AuthRepository interface {
	Login(ctx context.Context, username, password string) (data model.User, err error)
	UpdateLastLogin(ctx context.Context, userID uint, at time.Time) error
	ChangeMyPassword(ctx context.Context, id uint, password string) error
	GetMe(ctx context.Context, id uint) (data model.User, err error)
	GetRolePermissions(ctx context.Context, userID uint) ([]string, []string, error)
	HasUserPermission(ctx context.Context, userID uint, code string) (bool, error)
	ListUserMenus(ctx context.Context, userID uint) ([]model.Menu, error)
	ListAllMenus(ctx context.Context) ([]model.Menu, error)
	GetMenusByIDs(ctx context.Context, ids []uint) ([]model.Menu, error)
	WithTx(tx *gorm.DB) AuthRepository // 提供一个事务操作接口
}

// AuthRepository 实现了 AuthRepository 接口
type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// Login 登录
func (r *authRepository) Login(ctx context.Context, username, password string) (data model.User, err error) {
	// 统一错误提示，防止用户枚举
	const authErrorMsg = "账号或密码错误，请重新输入"
	if err = r.db.WithContext(ctx).Where("username = ?", username).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New(authErrorMsg)
		}
	}
	// 非活跃用户禁止登录
	if data.IsActive != nil && !*data.IsActive {
		return data, errors.New(authErrorMsg)
	}
	// 密码校验
	if err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(password)); err != nil {
		return data, errors.New(authErrorMsg)
	}
	return
}

// UpdateLastLogin 登录成功后更新最后登录时间
func (r *authRepository) UpdateLastLogin(ctx context.Context, userID uint, at time.Time) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("last_login", at).Error
}

// ChangeMyPassword 修改密码
func (r *authRepository) ChangeMyPassword(ctx context.Context, id uint, password string) error {
	// 修改密码
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Select("password").Update("password", password).Error; err != nil {
		return err
	}
	return nil
}

func (r *authRepository) GetMe(ctx context.Context, id uint) (data model.User, err error) {
	if err = r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (r *authRepository) GetRolePermissions(ctx context.Context, userID uint) (roles, permissions []string, err error) {
	if err = r.db.WithContext(ctx).
		Table("valyria_authn_role r").
		Joins("JOIN valyria_authn_user_role ur ON ur.role_id = r.id").
		Where("ur.user_id = ?", userID).
		Distinct("r.code").
		Pluck("r.code", &roles).Error; err != nil {
		return nil, nil, err
	}
	if err = r.db.WithContext(ctx).
		Table("valyria_authn_role_permission rp").
		Joins("JOIN valyria_authn_user_role ur ON ur.role_id = rp.role_id").
		Joins("JOIN valyria_authn_permission p ON p.id = rp.permission_id").
		Where("ur.user_id = ?", userID).
		Distinct("p.code").
		Pluck("p.code", &permissions).Error; err != nil {
		return nil, nil, err
	}
	return roles, permissions, nil
}

func (r *authRepository) HasUserPermission(ctx context.Context, userID uint, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("valyria_authn_role_permission rp").
		Joins("JOIN valyria_authn_user_role ur ON ur.role_id = rp.role_id").
		Joins("JOIN valyria_authn_permission p ON p.id = rp.permission_id").
		Where("ur.user_id = ? AND p.code = ?", userID, code).
		Count(&count).Error
	return count > 0, err
}

func (r *authRepository) ListUserMenus(ctx context.Context, userID uint) ([]model.Menu, error) {
	var data []model.Menu
	err := r.db.WithContext(ctx).
		Table("valyria_authn_menu m").
		Joins("JOIN valyria_authn_user_menu um ON um.menu_id = m.id").
		Where("um.user_id = ?", userID).
		Order("m.`rank` ASC, m.id ASC").
		Scan(&data).Error
	return data, err
}

func (r *authRepository) ListAllMenus(ctx context.Context) ([]model.Menu, error) {
	var data []model.Menu
	err := r.db.WithContext(ctx).
		Table("valyria_authn_menu").
		Scan(&data).Error
	return data, err
}

func (r *authRepository) GetMenusByIDs(ctx context.Context, ids []uint) ([]model.Menu, error) {
	var data []model.Menu
	if len(ids) == 0 {
		return data, nil
	}
	err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Order("`rank` ASC, id ASC").
		Find(&data).Error
	return data, err
}

// WithTx 返回一个绑定事务的 Repository
func (r *authRepository) WithTx(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}
