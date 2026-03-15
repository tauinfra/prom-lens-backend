package model

import (
	"time"
	"valyria-backend/internal/apps/authn/request"
)

type Permission struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	Code      string    `gorm:"size:64;unique;not null" json:"code"` // 例如 "user:add", "user:reset_password"
	Creator   string    `gorm:"type:varchar(64);not null;comment:创建人" json:"creator,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Permission) TableName() string {
	return "valyria_authn_permission"
}

func (p *Permission) ApplyRequest(req *request.UpdatePermissionRequest) {
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Code != nil {
		p.Code = *req.Code
	}
}
