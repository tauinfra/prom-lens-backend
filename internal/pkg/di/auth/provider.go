package auth

import (
	"valyria-backend/internal/core/middleware"

	"gorm.io/gorm"
)

type AuthProvider struct {
	User   *UserProvider
	Role   *RuleProvider
	Module *ModuleProvider
	// ... 其他资源
}

func NewAuthProvider(db *gorm.DB, jwtManager *middleware.JWTManager) *AuthProvider {
	return &AuthProvider{
		User:   NewUserProvider(db, jwtManager),
		Role:   NewRuleProvider(db),
		Module: NewModuleProvider(db),
	}
}
