package model

import (
	"time"

	"gorm.io/datatypes"
)

type Rule struct {
	ID          uint           `gorm:"type:bigint;primaryKey"                             json:"id"`
	Name        string         `gorm:"type:varchar(64);not null;uniqueIndex:uniq_prom_rule_group_name;comment:'告警名称'" json:"name,omitempty"`
	GroupID     int            `gorm:"not null;uniqueIndex:uniq_prom_rule_group_name;comment:'规则组'" json:"groupID,omitempty"`
	Summary     string         `gorm:"type:varchar(255);not null;comment:'告警描述'"       json:"summary,omitempty"`
	Description string         `gorm:"type:longtext;not null;comment:'告警详情'"           json:"description,omitempty"`
	Expr        string         `gorm:"type:longtext;not null;comment:'告警规则'"           json:"expr,omitempty"`
	For         string         `gorm:"type:varchar(4);not null;comment:'持续时间'"         json:"for,omitempty"`
	Labels            datatypes.JSON `gorm:"type:json;not null;comment:'规则标签'"                         json:"labels,omitempty"`
	ExtraAnnotations  datatypes.JSON `gorm:"column:extra_annotations;type:json;not null;comment:'扩展 annotations'" json:"extraAnnotations,omitempty"`
	Status            *bool          `gorm:"type:tinyint(1);default:1;comment:'规则状态'"                  json:"status"` // 启用状态
	CreatedAt   *time.Time     `gorm:"datetime(3);comment:'创建时间'"                      json:"createdAt"`
	UpdatedAt   *time.Time     `gorm:"datetime(3);comment:'更新时间'"                      json:"updatedAt"`
}

func (Rule) TableName() string {
	return "prom_lens_prometheus_rule"
}
