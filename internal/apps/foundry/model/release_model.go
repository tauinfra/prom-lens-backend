package model

import (
	"time"
	"valyria-backend/internal/apps/auth/model"
)

type Release struct {
	ID                 int         `gorm:"type:bigint;primaryKey"                               json:"id,omitempty"`
	Pipeline           *Pipeline   `gorm:"foreignKey:PipelineID"                            json:"pipeline,omitempty"`
	PipelineID         int         `gorm:"type:int;not null;comment:'流水线'" json:"pipelineID,omitempty"`                 // 对应 Pipeline
	TaskID             string      `gorm:"type:varchar(64);not null;unique;comment:'TaskID'"   json:"taskID,omitempty"` // PipelineRun 生成的唯一 ID
	GitRef             string      `gorm:"type:varchar(128);not null;comment:'Git 分支或 Tag'" json:"gitRef,omitempty"`
	GitRefType         string      `gorm:"type:enum('branch','tag');not null;comment:'Git Ref 类型'" json:"gitRefType,omitempty"`
	GitCommit          string      `gorm:"type:char(40);not null;comment:'Git Commit SHA'" json:"gitCommit,omitempty"`
	ImageRegistry      string      `gorm:"type:varchar(128);not null;comment:'镜像仓库'" json:"imageRegistry,omitempty"`
	ImageName          string      `gorm:"type:varchar(128);not null;comment:'镜像名称'" json:"imageName,omitempty"`
	ImageTag           string      `gorm:"type:varchar(128);not null;comment:'镜像标签'" json:"imageTag,omitempty"`
	Status             string      `gorm:"type:enum('pending','progressing','success','failed');not null;default:'pending';comment:'发布状态'" json:"status,omitempty"` // 审核状态: pending（待审核）、approved（审核通过）、rejected（审核拒绝）
	DeployedBy         string      `gorm:"type:varchar(64);not null;comment:'发布人员'"         json:"deployedBy,omitempty" validate:"required"`                        //
	DeployedByUsername *model.User `gorm:"foreignKey:DeployedBy;references:Username"       json:"deployedByUsername,omitempty"`
	Description        string      `gorm:"type:varchar(255);comment:'发布描述'"                  json:"description" validate:"required"` // 发布描述
	StartedAt          *time.Time  `gorm:"type:datetime(3);comment:'开始时间'"                   json:"startedAt,omitempty"`             // 开始时间
	FinishedAt         *time.Time  `gorm:"type:datetime(3);comment:'结束时间'"                   json:"finishedAt,omitempty"`
	CreatedAt          *time.Time  `gorm:"type:datetime(3);comment:'创建时间'"                   json:"createdAt"`
	UpdatedAt          *time.Time  `gorm:"type:datetime(3);comment:'更新时间'"                   json:"updatedAt"`
}

func (Release) TableName() string {
	return "valyria_foundry_release"
}
