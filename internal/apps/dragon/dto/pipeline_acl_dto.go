package dto

import "time"

type PipelineACLDTO struct {
	ID           uint      `json:"id"`
	PipelineID   uint      `json:"pipelineID"`
	PipelineName string    `json:"pipelineName"`
	UserID       uint      `json:"userID"`
	UserName     string    `json:"userName"`
	Action       string    `json:"action"`
	Creator      string    `json:"creator"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
