package authn

import (
	"valyria-backend/internal/core/middleware"

	"gorm.io/gorm"
)

type AuthnProvider struct {
	Auth       *AuthProvider // 登录
	User       *UserProvider
	Role       *RoleProvider
	Menu       *MenuProvider
	Permission *PermissionProvider
	// ... 其他资源
}

func NewAuthnProvider(db *gorm.DB, jwtManager *middleware.JWTManager) *AuthnProvider {
	return &AuthnProvider{
		Auth:       NewAuthProvider(db, jwtManager),
		User:       NewUserProvider(db),
		Role:       NewRoleProvider(db), // 后端 API 操作权限角色
		Menu:       NewMenuProvider(db),
		Permission: NewPermissionProvider(db),
	}
}
