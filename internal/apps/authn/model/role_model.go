package model

import (
	"time"
	"valyria-backend/internal/apps/authn/request"
)

type Role struct {
	ID          uint      `gorm:"type:bigint;primaryKey"                                 json:"id,omitempty"`
	Name        string    `gorm:"type:varchar(32);not null;unique;comment:'角色名称'"     json:"name,omitempty"`
	Code        string    `gorm:"type:varchar(64);not null;unique;comment:'角色编码'"     json:"code,omitempty"`
	Description string    `gorm:"type:varchar(255);comment:'角色描述'"                 json:"description,omitempty"`
	Creator     string    `gorm:"type:varchar(64);not null;comment:创建人" json:"creator,omitempty"`
	CreatedAt   time.Time `gorm:"datetime(3);comment:'创建时间'"                          json:"createdAt,omitempty"`
	UpdatedAt   time.Time `gorm:"datetime(3);comment:'更新时间'"                          json:"updatedAt,omitempty"`
}

func (Role) TableName() string {
	return "valyria_authn_role"
}

type RolePermission struct {
	RoleID       uint      `gorm:"primaryKey" json:"roleID"`
	PermissionID uint      `gorm:"primaryKey" json:"permissionID"`
	Creator      string    `gorm:"type:varchar(64);not null;comment:创建人" json:"creator,omitempty"`
	CreatedAt    time.Time `gorm:"datetime(3);comment:'创建时间'"                    json:"createdAt,omitempty"`
	UpdatedAt    time.Time `gorm:"datetime(3);comment:'更新时间'"                    json:"updatedAt,omitempty"`
}

func (RolePermission) TableName() string {
	return "valyria_authn_role_permission"
}

func (p *Role) ApplyRequest(req *request.UpdateRoleRequest) {
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Code != nil {
		p.Code = *req.Code
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
}
