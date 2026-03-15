package request

import (
	domainpermission "valyria-backend/internal/apps/kubernetes/domain/permission"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

const permissionRoleTag = "permission_role"

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation(permissionRoleTag, func(fl validator.FieldLevel) bool {
			return domainpermission.IsValidRole(fl.Field().String())
		})
	}
}
