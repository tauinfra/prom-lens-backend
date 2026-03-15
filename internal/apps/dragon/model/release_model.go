package model

import (
	"time"
)

type Release struct {
	ID             uint       `gorm:"type:bigint;primaryKey"                               json:"id,omitempty"`
	Pipeline       *Pipeline  `gorm:"foreignKey:PipelineID"                            json:"pipeline,omitempty"`
	PipelineID     uint       `gorm:"type:bigint;not null;comment:'流水线 ID'" json:"pipelineID,omitempty"`
	TaskID         string     `gorm:"type:varchar(64);not null;unique;comment:'TaskID'"   json:"taskID,omitempty"`
	GitRef         string     `gorm:"type:varchar(128);not null;comment:'Git 分支或 Tag'" json:"gitRef,omitempty"`
	GitRefType     string     `gorm:"type:enum('branch','tag');not null;comment:'Git Ref 类型'" json:"gitRefType,omitempty"`
	GitCommit      string     `gorm:"type:char(40);not null;comment:'Git Commit SHA'" json:"gitCommit,omitempty"`
	ImageRegistry  string     `gorm:"type:varchar(128);not null;comment:'镜像仓库'" json:"imageRegistry,omitempty"`
	ImageName      string     `gorm:"type:varchar(128);not null;comment:'镜像名称'" json:"imageName,omitempty"`
	ImageTag       string     `gorm:"type:varchar(128);not null;comment:'镜像标签'" json:"imageTag,omitempty"`
	ReleaseStatus  string     `gorm:"type:enum('pending','progressing','success','failed');not null;default:'pending';comment:'发布状态'" json:"releaseStatus,omitempty"`
	ApprovalStatus string     `gorm:"type:enum('none','pending','approved','rejected');not null;default:'none';comment:'审批状态'" json:"approvalStatus,omitempty"` // 审批状态: none（无审批）、pending（待审核）、approved（审核通过）、rejected（审核拒绝）
	Description         string     `gorm:"type:varchar(255);comment:'发布描述'"                  json:"description"`                                                     // 发布描述
	TargetReleaseID     *uint      `gorm:"type:bigint;default:null;comment:'回滚目标 release id'" json:"targetReleaseID,omitempty"`
	Creator             string     `gorm:"type:varchar(64);not null;comment:创建人" json:"creator,omitempty"`
	StartedAt      *time.Time `gorm:"type:datetime(3);comment:'开始时间'"                   json:"startedAt,omitempty"` // 开始时间
	FinishedAt     *time.Time `gorm:"type:datetime(3);comment:'结束时间'"                   json:"finishedAt,omitempty"`
	CreatedAt      time.Time  `gorm:"type:datetime(3);comment:'创建时间'"                   json:"createdAt,omitempty"`
	UpdatedAt      time.Time  `gorm:"type:datetime(3);comment:'更新时间'"                   json:"updatedAt,omitempty"`
}

func (Release) TableName() string {
	return "valyria_dragon_release"
}
