package model

import (
	"gorm.io/datatypes"
	"time"
)

type AuditLog struct {
	ID         int            `gorm:"type:bigint;primaryKey"   json:"id,omitempty"`
	Username   string         `gorm:"type:varchar(32);not null;comment:'登录账号'"  json:"username,omitempty" validate:"required"`
	UrlPath    string         `gorm:"type:varchar(255);not null;comment:'URL地址" json:"urlPath,omitempty" validate:"required"`
	Method     string         `gorm:"type:varchar(32);not null;comment:'请求方法"   json:"method,omitempty" validate:"required"`
	IPAddress  string         `gorm:"type:varchar(64);not null;comment:'登录地址'" json:"IPAddress,omitempty" validate:"required"`
	Agent      string         `gorm:"type:varchar(255);not null;comment:'客户端'"  json:"agent,omitempty" validate:"required"`
	StatusCode int            `gorm:"type:int;not null;comment:'状态码'"          json:"statusCode,omitempty" validate:"required"` // 状态码
	Params     datatypes.JSON `gorm:"type:json"                  json:"params,omitempty"`
	Response   datatypes.JSON `gorm:"type:json;not null"         json:"response,omitempty"`
	CreatedAt  time.Time      `gorm:"datetime(3);comment:'创建时间'"             json:"createdAt,omitempty"`
}

func (AuditLog) TableName() string {
	return "valyria_audit_log"
}
