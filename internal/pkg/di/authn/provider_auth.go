package authn

import (
	"prom-lens-backend/internal/apps/authn/controller"
	"prom-lens-backend/internal/apps/authn/repository"
	"prom-lens-backend/internal/apps/authn/service"
	"prom-lens-backend/internal/core/middleware"

	"gorm.io/gorm"
)

type AuthProvider struct {
	Repo       repository.AuthRepository
	Service    service.AuthManager
	Controller *controller.AuthController
}

func NewAuthProvider(db *gorm.DB, jwtManager *middleware.JWTManager) *AuthProvider {
	repo := repository.NewAuthRepository(db)
	manager := service.NewAuthManager(db, repo)
	control := controller.NewAuthController(manager, jwtManager)

	return &AuthProvider{
		Repo:       repo,
		Service:    manager,
		Controller: control,
	}
}
