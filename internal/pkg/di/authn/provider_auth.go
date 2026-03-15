package authn

import (
	"valyria-backend/internal/apps/authn/controller"
	"valyria-backend/internal/apps/authn/repository"
	"valyria-backend/internal/apps/authn/service"
	"valyria-backend/internal/core/middleware"

	"gorm.io/gorm"
)

type AuthProvider struct {
	Repo       *repository.AuthRepository
	Service    *service.AuthManager
	Controller *controller.AuthController
}

func NewAuthProvider(db *gorm.DB, jwtManager *middleware.JWTManager) *AuthProvider {
	// 初始化用户模块的依赖
	repo := repository.NewAuthRepository(db)
	manager := service.NewAuthManager(db, repo)
	control := controller.NewAuthController(manager, jwtManager)

	return &AuthProvider{
		Repo:       &repo,
		Service:    &manager,
		Controller: control,
	}
}
