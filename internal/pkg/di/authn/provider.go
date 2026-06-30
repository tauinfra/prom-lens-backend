package authn

import (
	"prom-lens-backend/internal/core/middleware"

	"gorm.io/gorm"
)

type AuthnProvider struct {
	Auth *AuthProvider
	User *UserProvider
}

func NewAuthnProvider(db *gorm.DB, jwtManager *middleware.JWTManager) *AuthnProvider {
	return &AuthnProvider{
		Auth: NewAuthProvider(db, jwtManager),
		User: NewUserProvider(db),
	}
}
