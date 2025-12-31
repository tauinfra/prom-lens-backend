package model

import (
	"time"
)

type Environment struct {
	ID           int        `gorm:"type:bigint;primaryKey"                                 json:"id,omitempty"`
	Name         string     `gorm:"type:varchar(32);not null;uniqueIndex:uk_project_env;comment:'环境名称'" json:"name"`
	ProjectID    int        `gorm:"type:bigint;not null;uniqueIndex:uk_project_env" json:"projectID"`
	Project      Project    `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	CredArgoCDID int        `json:"credArgoCDID,omitempty"`
	CredArgoCD   CredArgoCD `gorm:"foreignKey:CredArgoCDID" json:"credArgoCD,omitempty"`
	CredHarborID int        `json:"credHarborID,omitempty"`
	CredHarbor   CredHarbor `gorm:"foreignKey:CredHarborID;not null;comment:'镜像仓库'" json:"credHarbor,omitempty"` // 使用指针，空字段不显示
	Description  string     `gorm:"type:text;comment:'环境描述'"                            json:"description,omitempty"`
	CreatedAt    *time.Time `gorm:"datetime(3);comment:'创建时间'"                          json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `gorm:"datetime(3);comment:'更新时间'"                          json:"updatedAt,omitempty"`
}

func (Environment) TableName() string {
	return "valyria_foundry_environment"
}
