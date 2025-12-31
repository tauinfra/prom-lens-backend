package model

import (
	"time"
)

type User struct {
	ID          int        `gorm:"type:bigint;primaryKey"                                 json:"id,omitempty"`
	Username    string     `gorm:"type:varchar(32);not null;unique;comment:'用户名称'"     json:"username,omitempty" validate:"required"`
	Password    string     `gorm:"type:varchar(128);comment:'用户密码'"                    json:"password,omitempty"` // 启动 LDAP 认证, 密码可以为空
	Nickname    string     `gorm:"type:varchar(32);comment:'中文名称'"                     json:"nickname,omitempty"`
	Email       string     `gorm:"type:varchar(64);not null;unique;comment:'邮箱地址'"     json:"email,omitempty" validate:"email"`
	Phone       string     `gorm:"type:varchar(11);default:null;unique;comment:'手机号码'" json:"phone,omitempty"`
	IsActive    *bool      `gorm:"type:tinyint(1);default:1;comment:'活跃状态'"            json:"isActive,omitempty"`   // 是否允许用户登录
	IsSuperuser *bool      `gorm:"type:tinyint(1);default:0;comment:'超级管理员'"          json:"isSuperuser,omitempty"` // 是否为超级用户
	IsLdap      *bool      `gorm:"type:tinyint(1);default:0;comment:'LDAP 认证'"          json:"isLdap,omitempty"`    // 是否为 LDAP 认证
	Modules     []Module   `gorm:"many2many:valyria_auth_user_module" json:"modules,omitempty"`
	ModuleIDs   []int      `gorm:"-" json:"moduleIds,omitempty"` // 只用于前端输入（绑定用）保存仍使用 Modules []Module
	Roles       []Role     `gorm:"many2many:valyria_auth_user_role" json:"roles,omitempty"`
	RoleIDs     []int      `gorm:"-" json:"roleIds,omitempty"` // 只用于前端输入（绑定用）保存仍使用 Roles []Role
	CreatedAt   *time.Time `gorm:"datetime(3);comment:'创建时间'"                          json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `gorm:"datetime(3);comment:'更新时间'"                          json:"updatedAt,omitempty"`
}

type UserModule struct {
	UserID    int        `gorm:"primaryKey" json:"userID"`
	ModuleID  int        `gorm:"primaryKey" json:"moduleID"`
	CreatedAt *time.Time `gorm:"datetime(3);comment:'创建时间'"                    json:"createdAt,omitempty"`
	UpdatedAt *time.Time `gorm:"datetime(3);comment:'更新时间'"                    json:"updatedAt,omitempty"`
}

type UserRole struct {
	UserID    int        `gorm:"primaryKey" json:"userID"`
	RoleID    int        `gorm:"primaryKey" json:"roleID"`
	CreatedAt *time.Time `gorm:"datetime(3);comment:'创建时间'"                    json:"createdAt,omitempty"`
	UpdatedAt *time.Time `gorm:"datetime(3);comment:'更新时间'"                    json:"updatedAt,omitempty"`
}

func (User) TableName() string {
	return "valyria_auth_user"
}

func (UserModule) TableName() string {
	return "valyria_auth_user_module"
}

func (UserRole) TableName() string { return "valyria_auth_user_role" }
