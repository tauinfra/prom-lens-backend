package model

import (
	"time"

	"prom-lens-backend/internal/apps/authn/request"
)

type User struct {
	ID          int        `gorm:"type:bigint;primaryKey"                                 json:"id,omitempty"`
	Username    string     `gorm:"type:varchar(32);not null;unique;comment:'用户名称'"     json:"username,omitempty"`
	Password    string     `gorm:"type:varchar(128);comment:'用户密码'"                    json:"-"`
	Nickname    string     `gorm:"type:varchar(32);comment:'中文名称'"                     json:"nickname,omitempty"`
	Email       string     `gorm:"type:varchar(64);not null;unique;comment:'邮箱地址'"     json:"email,omitempty"`
	Phone       string     `gorm:"type:varchar(11);default:null;unique;comment:'手机号码'" json:"phone,omitempty"`
	IsActive    *bool      `gorm:"type:tinyint(1);default:1;comment:'活跃状态'"            json:"isActive,omitempty"`
	IsSuperuser *bool      `gorm:"type:tinyint(1);default:0;comment:'超级管理员'"          json:"isSuperuser,omitempty"`
	IsLdap      *bool      `gorm:"type:tinyint(1);default:0;comment:'LDAP 认证'"          json:"isLdap,omitempty"`
	DN          *string    `gorm:"type:varchar(255);uniqueIndex;comment:'LDAP DN'" json:"dn,omitempty"`
	Creator     string     `gorm:"type:varchar(64);not null;comment:创建人" json:"creator,omitempty"`
	LastLogin   *time.Time `gorm:"type:datetime(3);comment:'最后登录时间'" json:"lastLogin,omitempty"`
	CreatedAt   time.Time  `gorm:"datetime(3);comment:'创建时间'"                          json:"createdAt,omitempty"`
	UpdatedAt   time.Time  `gorm:"datetime(3);comment:'更新时间'"                          json:"updatedAt,omitempty"`
}

func (User) TableName() string {
	return "prom_lens_authn_user"
}

func (u *User) ApplyRequest(req *request.UpdateUserRequest) {
	if req.Username != nil {
		u.Username = *req.Username
	}
	if req.Nickname != nil {
		u.Nickname = *req.Nickname
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.Phone != nil {
		u.Phone = *req.Phone
	}
	if req.IsActive != nil {
		u.IsActive = req.IsActive
	}
	if req.IsSuperuser != nil {
		u.IsSuperuser = req.IsSuperuser
	}
	if req.IsLdap != nil {
		u.IsLdap = req.IsLdap
	}
	if req.DN != nil {
		u.DN = req.DN
	}
}
