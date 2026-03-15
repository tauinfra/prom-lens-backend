package request

// CreateReviewRequest 创建审批记录请求（仅允许传 releaseID，申请人由服务端从 ctx 写入）
type CreateReviewRequest struct {
	ReleaseID uint `json:"releaseID" binding:"required"`
}

// UpdateReviewRequest 更新审批请求（仅允许传 status、comment，审批人由服务端从 ctx 写入）
type UpdateReviewRequest struct {
	Status  string `json:"status" binding:"required,oneof=approved rejected"`
	Comment string `json:"comment"`
}
