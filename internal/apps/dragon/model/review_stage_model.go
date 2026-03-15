package model

import (
	"time"
	"valyria-backend/internal/apps/authn/model"
)

// ReviewStageConfig 全局审批阶段配置
type ReviewStageConfig struct {
	ID        uint        `gorm:"type:bigint;primaryKey" json:"id,omitempty"`
	Step      int         `gorm:"type:tinyint(2);not null;uniqueIndex:uniq_review_stage_step" json:"step,omitempty"`
	RoleID    uint        `gorm:"not null" json:"roleID,omitempty"`
	Role      *model.Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Creator   string      `gorm:"type:varchar(64);not null;comment:创建人" json:"creator,omitempty"`
	CreatedAt time.Time   `gorm:"datetime(3);comment:'创建时间'" json:"createdAt,omitempty"`
	UpdatedAt time.Time   `gorm:"datetime(3);comment:'更新时间'" json:"updatedAt,omitempty"`
}

func (ReviewStageConfig) TableName() string {
	return "valyria_dragon_review_stage"
}
