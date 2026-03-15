package model

import (
	"time"

	"gorm.io/datatypes"
)

// TargetGroup 组
type TargetGroup struct {
	ID          int            `gorm:"type:bigint;primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(64);not null;unique;comment:'Target 分组'" json:"name,omitempty"`
	Description string         `gorm:"type:varchar(64);comment:'Target 分组描述'" json:"description,omitempty"`
	Labels      datatypes.JSON `gorm:"type:json;not null;comment:'分组标签'" json:"labels"`
	CreatedAt   time.Time      `gorm:"datetime(3);comment:'创建时间'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time      `gorm:"datetime(3);comment:'更新时间'" json:"updatedAt,omitempty"`
}

func (TargetGroup) TableName() string {
	return "valyria_prometheus_target_group"
}
