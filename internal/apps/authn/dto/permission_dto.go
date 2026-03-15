package dto

import (
	"time"
)

type PermissionDTO struct {
	ID        uint      `gorm:"column:id" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Code      string    `gorm:"column:code" json:"code"`
	Creator   string    `gorm:"column:creator" json:"creator"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
