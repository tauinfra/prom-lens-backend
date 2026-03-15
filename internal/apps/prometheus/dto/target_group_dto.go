package dto

import (
	"time"

	"gorm.io/datatypes"
)

type TargetGroupDTO struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Labels      datatypes.JSON `json:"labels"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}
