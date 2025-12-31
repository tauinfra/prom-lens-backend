package model

import "time"

type Project struct {
	ID           int           `gorm:"type:bigint;primaryKey"                                 json:"id,omitempty"`
	Name         string        `gorm:"type:varchar(32);not null;unique;comment:'项目名称'"     json:"name,omitempty" validate:"required"`
	Environments []Environment `gorm:"foreignKey:ProjectID" json:"environments,omitempty"`
	Description  string        `gorm:"type:text;comment:'项目描述'"                            json:"description,omitempty"`
	CreatedAt    *time.Time    `gorm:"datetime(3);comment:'创建时间'"                          json:"createdAt,omitempty"`
	UpdatedAt    *time.Time    `gorm:"datetime(3);comment:'更新时间'"                          json:"updatedAt,omitempty"`
}

func (Project) TableName() string {
	return "valyria_foundry_project"
}
