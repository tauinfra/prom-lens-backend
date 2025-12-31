package model

import "time"

type Application struct {
	ID          int        `gorm:"type:bigint;primaryKey"              json:"id,omitempty"`
	Name        string     `gorm:"type:varchar(64);not null;comment:'应用名称'"           json:"name,omitempty"` // 应用名称
	Project     Project    `gorm:"foreignKey:ProjectID;not null" json:"project,omitempty"`
	ProjectID   int        `json:"projectID,omitempty"`                                                             // 应用 ID(写入)
	Description string     `gorm:"type:varchar(128);comment:'应用描述'"                   json:"description,omitempty"` // 应用描述
	CreatedAt   *time.Time `gorm:"datetime(3);comment:'创建时间'"       json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `gorm:"datetime(3);comment:'更新时间'"       json:"updatedAt,omitempty"`
}

func (Application) TableName() string {
	return "valyria_foundry_application"
}
