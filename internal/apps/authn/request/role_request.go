package request

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name"`
	Code        *string `json:"code"`
	Description *string `json:"description"`
}

type UpdateRolePermissionsRequest struct {
	Permissions []uint `json:"permissions"` // 权限 ID 列表
}

type UpdateUserRolesRequest struct {
	Roles []uint `json:"roles"` // 角色 ID 列表
}

type UpdateUserMenusRequest struct {
	Menus []uint `json:"menus"` // 菜单 ID 列表
}
