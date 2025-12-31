package model

import "time"

type Pipeline struct {
	ID              int         `gorm:"type:bigint;primaryKey"                               json:"id,omitempty"`
	Name            string      `gorm:"type:varchar(256);not null;comment:'流水线名称'"          json:"name,omitempty"`
	ApplicationID   int         `json:"applicationID,omitempty"`
	Application     Application `gorm:"foreignKey:ApplicationID"                            json:"application,omitempty"`
	EnvironmentID   int         `json:"environmentID,omitempty"`
	Environment     Environment `gorm:"foreignKey:EnvironmentID" json:"environment,omitempty"`
	CredGitlabID    int         `json:"credGitlabID,omitempty"`
	CredGitlab      CredGitlab  `gorm:"foreignKey:CredGitlabID;comment:'GitLab 实例'" json:"credGitlab,omitempty"`
	GitlabGroupID   int         `gorm:"type:bigint;not null;comment:'GitLab 组织'" json:"gitlabGroupID,omitempty"`
	GitlabProjectID int         `gorm:"type:bigint;not null;comment:'GitLab 仓库'" json:"gitlabProjectID,omitempty"`
	IsApproval      *bool       `gorm:"type:bool;not null;default:false;comment:'是否启用审批'" json:"isApproval"`
	Template        string      `gorm:"type:varchar(256);not null;comment:'流水线模板'"          json:"template,omitempty"` // Pipeline 模板
	Script          string      `gorm:"type:text;comment:not null;'构建脚本'" json:"script,omitempty"`
	CreatedAt       *time.Time  `gorm:"type:datetime(3);comment:'创建时间'"                   json:"createdAt,omitempty"`
	UpdatedAt       *time.Time  `gorm:"type:datetime(3);comment:'更新时间'"                   json:"updatedAt,omitempty"`
}

func (Pipeline) TableName() string {
	return "valyria_foundry_pipeline"
}
