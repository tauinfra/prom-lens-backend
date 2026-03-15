package request

type CreatePipelineACLRequest struct {
	UserID uint   `json:"userID" binding:"required"`
	Action string `json:"action" binding:"required"`
}

type BatchCreatePipelineACLRequest struct {
	Users  []uint `json:"users" binding:"required"`
	Action string `json:"action" binding:"required"`
}
