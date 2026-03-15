package model

import (
	"time"
	"valyria-backend/internal/apps/dragon/request"
)

type Environment struct {
	ID          uint        `gorm:"primaryKey" json:"id,omitempty"`
	Name        string      `gorm:"type:varchar(32);not null;uniqueIndex:uk_project_env;comment:'环境名称'" json:"name" binding:"required"`
	ProjectID   uint        `gorm:"not null;uniqueIndex:uk_project_env" json:"projectID"`
	Project     *Project    `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;comment:'所属项目'" json:"project,omitempty" binding:"-"`
	HarborID    uint        `gorm:"comment:'Harbor 凭证 ID'" json:"harborID,omitempty"`
	Harbor      *Credential `gorm:"foreignKey:HarborID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:'Harbor 凭证'" json:"harbor,omitempty"`
	Description string      `gorm:"type:text;comment:'环境描述'" json:"description,omitempty"`
	Creator     string      `gorm:"type:varchar(64);not null;comment:创建人" json:"creator,omitempty"`
	CreatedAt   time.Time   `gorm:"autoCreateTime;comment:'创建时间'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime;comment:'更新时间'" json:"updatedAt,omitempty"`
}

func (Environment) TableName() string {
	return "valyria_dragon_environment"
}

func (u *Environment) ApplyRequest(req *request.UpdateEnvironmentRequest) {
	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.HarborID != nil {
		u.HarborID = *req.HarborID
	}
	if req.Description != nil {
		u.Description = *req.Description
	}
}
