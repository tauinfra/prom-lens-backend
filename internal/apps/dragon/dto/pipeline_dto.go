package dto

import "time"

type PipelineDTO struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	EnvironmentID   uint      `json:"environmentID"`
	EnvironmentName string    `json:"environmentName"`
	ProjectID       uint      `json:"projectID"`
	ProjectName     string    `json:"projectName"`
	HarborBaseURL   string    `json:"harborBaseURL"`
	GitlabID        uint      `json:"gitlabID"`
	GitlabGroupID   uint      `json:"gitlabGroupID"`
	GitlabProjectID uint      `json:"gitlabProjectID"`
	GitlabBaseURL   string    `json:"gitlabBaseURL"`
	IsApproval      bool      `json:"isApproval"`
	Task            string    `json:"task"`
	Creator         string    `json:"creator"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
