package model

import "time"

type Record struct {
	ID        int        `gorm:"type:bigint;primaryKey" json:"id"`
	Name      string     `gorm:"type:varchar(64);not null;unique;comment:'记录名称'" json:"name"`
	GroupID   int        `gorm:"gorm:foreignKey:GroupID" json:"groupID,omitempty"`
	Expr      string     `gorm:"type:longtext;not null;comment:'记录规则'" json:"expr"`
	CreatedAt *time.Time `gorm:"datetime(3);comment:'创建时间'" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"datetime(3);comment:'更新时间'" json:"updatedAt"`
}

func (Record) TableName() string {
	return "valyria_prometheus_record"
}
