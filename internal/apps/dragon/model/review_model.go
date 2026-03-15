package model

import (
	"time"
)

type Review struct {
	ID           uint       `gorm:"type:bigint;primaryKey" json:"id,omitempty"`
	ApplicantID  uint       `json:"applicantID,omitempty"`
	ReviewerID   uint       `json:"reviewerID,omitempty"`
	Release      *Release   `gorm:"foreignKey:ReleaseID" json:"release,omitempty"`                                                                  // 一对多外键(读权限，查询 Environment 时需要)
	ReleaseID    uint       `gorm:"type:bigint;not null;comment:'发布任务'" json:"releaseID,omitempty"`                                                 // 一对多外键(写权限，迁移时需要引用外键存在的关联表)
	Status       string     `gorm:"type:enum('pending', 'approved', 'rejected');not null;default:'pending';comment:'审核状态'" json:"status,omitempty"` // 审核状态: pending（待审核）、approved（审核通过）、rejected（审核拒绝）
	ApprovalRole string     `gorm:"type:varchar(32);not null;comment:'当前步骤审批角色'" json:"approvalRole"`
	Step         int        `gorm:"type:tinyint(2);default:1;not null;comment:'当前审核步骤'" json:"step,omitempty"`
	MaxStep      int        `gorm:"type:tinyint(2);default:2;not null;comment:'审核总步骤数'" json:"maxStep,omitempty"`
	Comment      string     `gorm:"type:text;comment:'审核评论'" json:"comment,omitempty"`
	CreatedAt    *time.Time `gorm:"datetime(3);comment:'创建时间'"       json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `gorm:"datetime(3);comment:'更新时间'"       json:"updatedAt,omitempty"`
}

func (Review) TableName() string {
	return "valyria_dragon_review"
}
