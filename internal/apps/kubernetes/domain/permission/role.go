package permission

import "errors"

const (
	RoleViewer    = "val:viewer"
	RoleDeveloper = "val:developer"
	RoleOperator  = "val:operator"
	RoleAdmin     = "val:admin"
)

var (
	ValidRoles     = []string{RoleViewer, RoleDeveloper, RoleOperator, RoleAdmin}
	ErrInvalidRole = errors.New("角色必须是 val:viewer、val:developer、val:operator、val:admin 之一")
)

// IsValidRole 校验角色是否为允许的值
func IsValidRole(role string) bool {
	for _, r := range ValidRoles {
		if r == role {
			return true
		}
	}
	return false
}
