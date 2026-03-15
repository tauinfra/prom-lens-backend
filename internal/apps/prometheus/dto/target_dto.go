package dto

import (
	"time"

	"gorm.io/datatypes"
)

type TargetDTO struct {
	ID        uint           `json:"id"`
	GroupID   int            `json:"groupID"`
	IPAddress string         `json:"ipAddress"`
	Port      int            `json:"port"`
	Labels    datatypes.JSON `json:"labels"`
	Enabled   *bool          `json:"enabled"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}
