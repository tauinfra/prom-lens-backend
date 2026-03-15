package authn

import (
	"valyria-backend/internal/apps/authn/controller"
	"valyria-backend/internal/apps/authn/repository"
	"valyria-backend/internal/apps/authn/service"

	"gorm.io/gorm"
)

type UserProvider struct {
	Repo       *repository.UserRepository
	Service    *service.UserService
	Controller *controller.UserController
}

func NewUserProvider(db *gorm.DB) *UserProvider {
	// 初始化用户模块的依赖
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(db, userRepo)
	userController := controller.NewUserController(userService)

	return &UserProvider{
		Repo:       &userRepo,
		Service:    &userService,
		Controller: userController,
	}
}
