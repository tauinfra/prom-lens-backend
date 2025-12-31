package auth

import (
	"valyria-backend/internal/apps/auth/controller"
	"valyria-backend/internal/apps/auth/repository"
	"valyria-backend/internal/apps/auth/service"
	"valyria-backend/internal/core/middleware"

	"gorm.io/gorm"
)

type UserProvider struct {
	Repo       *repository.UserRepository
	Service    *service.UserService
	Controller *controller.UserController
}

func NewUserProvider(db *gorm.DB, jwtManager *middleware.JWTManager) *UserProvider {
	// 初始化用户模块的依赖
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(db, userRepo)
	userController := controller.NewUserController(userService, jwtManager)

	return &UserProvider{
		Repo:       &userRepo,
		Service:    &userService,
		Controller: userController,
	}
}
