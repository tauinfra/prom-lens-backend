package dto

import "time"

type EnvironmentDTO struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	ProjectID     uint      `json:"projectID"`
	ProjectName   string    `json:"projectName"`
	HarborID      uint      `json:"harborID"`
	HarborBaseURL string    `json:"harborBaseURL"`
	Description   string    `json:"description"`
	Creator       string    `json:"creator"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
