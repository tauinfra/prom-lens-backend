package model

import (
	"time"

	"gorm.io/datatypes"
)

type Target struct {
	ID        uint           `gorm:"type:bigint;primaryKey"                             json:"id"`
	GroupID   int            `gorm:"gorm:foreignKey:GroupID;comment:'Target 分组';uniqueIndex:uniq_prom_target_addr" json:"groupID,omitempty"`
	IPAddress string         `gorm:"type:varchar(255);not null;uniqueIndex:uniq_prom_target_addr;comment:'IP 地址'" json:"ip,omitempty"`
	Port      int            `gorm:"type:bigint;not null;uniqueIndex:uniq_prom_target_addr;comment:'端口'" json:"port,omitempty"`
	Labels    datatypes.JSON `gorm:"type:json;not null;comment:'实例标签'"               json:"labels,omitempty"`
	Enabled   *bool          `gorm:"type:tinyint(1);default:1;comment:'监控状态'"        json:"enabled,omitempty"` // 启用状态
	CreatedAt time.Time      `gorm:"datetime(3);comment:'创建时间'"                      json:"createdAt"`
	UpdatedAt time.Time      `gorm:"datetime(3);comment:'更新时间'"                      json:"updatedAt"`
}

func (Target) TableName() string {
	return "valyria_prometheus_target"
}
