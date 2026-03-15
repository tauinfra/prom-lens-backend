package request

type CreateClusterRequest struct {
	Name        string `json:"name" binding:"required"`
	Host        string `json:"host" binding:"required"`
	Token       string `json:"token" binding:"required"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type UpdateClusterRequest struct {
	Name        *string `json:"name"`
	Host        *string `json:"host"`
	Version     *string `json:"version"`
	Description *string `json:"description"`
	// Token 字段不包含在更新请求中，禁止修改
}

type UpdateClusterTokenRequest struct {
	Token string `json:"token" binding:"required"`
}
