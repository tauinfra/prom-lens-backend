package request

// CreateProjectRequest 创建项目请求（仅暴露允许客户端传入的字段）
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateProjectRequest 更新项目请求（PATCH 式，字段均为可选）
type UpdateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
