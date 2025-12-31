package model

import (
	"time"
)

type Cluster struct {
	ID          int        `gorm:"type:bigint;primaryKey;comment:'集群 ID'"            json:"id,omitempty"`
	Name        string     `gorm:"type:varchar(128);not null;unique;comment:'集群名称'" json:"name,omitempty" validate:"required"`
	Host        string     `gorm:"type:varchar(256);not null;unique;comment:'集群地址'" json:"host,omitempty" validate:"required"`
	Token       string     `gorm:"type:text;not null;comment:'集群凭证'"                json:"token,omitempty" validate:"required"`
	Version     string     `gorm:"type:varchar(64);comment:'集群版本'"                  json:"version,omitempty"`
	Description string     `gorm:"type:varchar(256);comment:'集群描述'"                 json:"description,omitempty"`
	CreatedAt   *time.Time `gorm:"type:datetime(3);comment:'创建时间'"                  json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `gorm:"type:datetime(3);comment:'更新时间'"                  json:"updatedAt,omitempty"`
}

func (Cluster) TableName() string {
	return "valyria_kubernetes_cluster"
}
