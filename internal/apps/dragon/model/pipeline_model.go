package model

import (
	"time"
	"valyria-backend/internal/apps/dragon/request"
)

type Pipeline struct {
	ID              uint         `gorm:"type:bigint;primaryKey"                               json:"id,omitempty"`
	Name            string       `gorm:"type:varchar(256);not null;comment:'发布名称'"          json:"name,omitempty"`
	EnvironmentID   uint         `json:"environmentID,omitempty"`
	Environment     *Environment `gorm:"foreignKey:EnvironmentID" json:"components,omitempty"`
	GitlabID        uint         `json:"gitlabID,omitempty"`
	Gitlab          *Credential  `gorm:"foreignKey:GitlabID;comment:'GitLab 实例'" json:"gitlab,omitempty"`
	GitlabGroupID   uint         `gorm:"type:bigint;not null;comment:'GitLab 组织'" json:"gitlabGroupID,omitempty"`
	GitlabProjectID uint         `gorm:"type:bigint;not null;comment:'GitLab 仓库'" json:"gitlabProjectID,omitempty"`
	IsApproval      bool         `gorm:"type:bool;not null;default:false;comment:'是否启用审批'" json:"isApproval"`
	Task            string       `gorm:"type:varchar(256);not null;comment:'发布任务'"          json:"task,omitempty"` // Pipeline 模板
	Script          string       `gorm:"type:text;comment:not null;'构建脚本'" json:"script,omitempty"`
	Creator         string       `gorm:"type:varchar(64);not null;comment:创建人" json:"creator,omitempty"`
	CreatedAt       time.Time    `gorm:"type:datetime(3);comment:'创建时间'"                   json:"createdAt,omitempty"`
	UpdatedAt       time.Time    `gorm:"type:datetime(3);comment:'更新时间'"                   json:"updatedAt,omitempty"`
}

type PipelineACL struct {
	ID         uint      `gorm:"type:bigint;primaryKey" json:"id,omitempty"`
	PipelineID uint      `gorm:"not null;index:idx_pipeline_user_action,unique" json:"pipelineID,omitempty"`
	UserID     uint      `gorm:"not null;index:idx_pipeline_user_action,unique" json:"userID,omitempty"`
	Action     string    `gorm:"type:varchar(16);not null;index:idx_pipeline_user_action,unique" json:"action,omitempty"`
	Creator    string    `gorm:"type:varchar(64);not null;comment:创建人" json:"creator,omitempty"`
	CreatedAt  time.Time `gorm:"datetime(3);comment:'创建时间'" json:"createdAt,omitempty"`
	UpdatedAt  time.Time `gorm:"datetime(3);comment:'更新时间'" json:"updatedAt,omitempty"`
}

func (PipelineACL) TableName() string {
	return "valyria_dragon_pipeline_acl"
}

func (Pipeline) TableName() string {
	return "valyria_dragon_pipeline"
}

func (p *Pipeline) ApplyUpdate(req *request.UpdatePipelineRequest) {
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.GitlabID != nil {
		p.GitlabID = *req.GitlabID
	}
	if req.GitlabGroupID != nil {
		p.GitlabGroupID = *req.GitlabGroupID
	}
	if req.GitlabProjectID != nil {
		p.GitlabProjectID = *req.GitlabProjectID
	}
	if req.IsApproval != nil {
		p.IsApproval = *req.IsApproval
	}
	if req.Task != nil {
		p.Task = *req.Task
	}
	if req.Script != nil {
		p.Script = *req.Script
	}
}
