package model

import "time"

type AuthLog struct {
	ID         int       `gorm:"type:bigint;primaryKey"   json:"id,omitempty"`
	Username   string    `gorm:"type:varchar(64);not null;comment:'登录账号'" json:"username,omitempty"`
	IPAddress  string    `gorm:"type:varchar(64);not null;comment:'登录地址'" json:"IPAddress,omitempty"`
	System     string    `gorm:"type:varchar(64);comment:'系统版本'"           json:"system,omitempty"`
	Agent      string    `gorm:"type:varchar(255);not null;comment:'客户端'"  json:"agent,omitempty"`
	StatusCode int       `gorm:"type:int;not null;comment:'状态码'"          json:"statusCode,omitempty"` // 状态码
	CreatedAt  time.Time `gorm:"datetime(3);comment:'创建时间'"             json:"createdAt,omitempty"`
}

func (AuthLog) TableName() string {
	return "valyria_audit_auth_log"
}
