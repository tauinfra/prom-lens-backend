package request

type CreatePermissionRequest struct {
	Name    string `json:"name" binding:"required"`
	Code    string `json:"code" binding:"required"`
	Creator string `json:"creator"`
}

type UpdatePermissionRequest struct {
	Name *string `json:"name" binding:"required"`
	Code *string `json:"code" binding:"required"`
}
