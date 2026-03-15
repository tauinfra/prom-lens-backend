package request

type CreateEnvironmentRequest struct {
	Name        string `json:"name" binding:"required"`
	ProjectID   uint   `json:"projectID"`
	HarborID    uint   `json:"harborID" binding:"required"`
	Description string `json:"description"`
	Creator     string `json:"creator"` // 建议由 controller 从 ctx 写入，不信任客户端
}

type UpdateEnvironmentRequest struct {
	Name        *string `json:"name"`
	HarborID    *uint   `json:"harborID"`
	Description *string `json:"description"`
	Script      *string `json:"script"`
}
