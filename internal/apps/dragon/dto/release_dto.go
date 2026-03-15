package dto

import "time"

type ReleaseDTO struct {
	ID              uint       `json:"id"`
	TaskID          string     `json:"taskID"`
	ProjectID       uint       `json:"projectID"`       // 项目 ID
	ProjectName     string     `json:"projectName"`     // 项目 Name
	EnvironmentID   uint       `json:"environmentID"`   // 环境 ID
	EnvironmentName string     `json:"environmentName"` // 环境 Name
	PipelineID      uint       `json:"pipelineID"`      // 流水线 ID
	PipelineName    string     `json:"pipelineName"`    // 流水线 Name
	GitRef          string     `json:"gitRef"`
	GitCommit       string     `json:"gitCommit"`
	GitRefType      string     `json:"gitRefType"`
	ImageRegistry   string     `json:"imageRegistry"`
	ImageName       string     `json:"imageName"`
	ImageTag        string     `json:"imageTag"`
	ReleaseStatus   string     `json:"releaseStatus"`
	ApprovalStatus   string     `json:"approvalStatus"` // 审批状态: none（无审批）、pending（待审核）、approved（审核通过）、rejected（审核拒绝）
	Description      string     `json:"description"`
	TargetReleaseID *uint      `json:"targetReleaseID,omitempty"` // 回滚目标 release id，非空表示本条为回滚记录
	Creator          string     `json:"creator"`
	StartedAt       *time.Time `json:"startedAt,omitempty"`  // 开始时间
	FinishedAt      *time.Time `json:"finishedAt,omitempty"` // 结束时间
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}
