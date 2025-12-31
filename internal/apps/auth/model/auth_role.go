package model

import "time"

type Role struct {
	ID          int        `gorm:"type:bigint;primaryKey"                                 json:"id,omitempty"`
	Name        string     `gorm:"type:varchar(32);not null;unique;comment:'角色名称'"     json:"name,omitempty" validate:"required"`
	Description string     `gorm:"type:varchar(255);comment:'角色描述'"                 json:"description,omitempty"`
	CreatedAt   *time.Time `gorm:"datetime(3);comment:'创建时间'"                          json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `gorm:"datetime(3);comment:'更新时间'"                          json:"updatedAt,omitempty"`
}

func (Role) TableName() string {
	return "valyria_auth_role"
}
