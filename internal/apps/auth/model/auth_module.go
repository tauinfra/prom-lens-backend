package model

import "time"

type Module struct {
	ID          int        `gorm:"type:bigint;primaryKey"                                 json:"id,omitempty"`
	Name        string     `gorm:"type:varchar(32);not null;unique;comment:'用户名称'"     json:"name,omitempty" validate:"required"`
	RoutePrefix string     `gorm:"type:varchar(128);comment:'路由前缀'"                    json:"routePrefix,omitempty"` // 启动 LDAP 认证, 密码可以为空
	Description string     `gorm:"type:varchar(32);comment:'模块描述'"                     json:"description,omitempty"`
	CreatedAt   *time.Time `gorm:"datetime(3);comment:'创建时间'"                          json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `gorm:"datetime(3);comment:'更新时间'"                          json:"updatedAt,omitempty"`
}

func (Module) TableName() string {
	return "valyria_auth_module"
}
