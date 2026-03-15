package dto

import "time"

// ReviewDTO 审批记录列表/详情响应
type ReviewDTO struct {
	ID           uint       `json:"id,omitempty"`
	ApplicantID  uint       `json:"applicantID,omitempty"`
	ReviewerID   uint       `json:"reviewerID,omitempty"`
	ReviewerName string     `json:"reviewerName,omitempty"` // 审批人姓名（来自 authn 用户）
	ReleaseID    uint       `json:"releaseID,omitempty"`
	Status       string     `json:"status,omitempty"`
	ApprovalRole string     `json:"approvalRole,omitempty"`
	Step         int        `json:"step,omitempty"`
	MaxStep      int        `json:"maxStep,omitempty"`
	Comment      string     `json:"comment,omitempty"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
}
