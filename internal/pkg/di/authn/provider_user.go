package authn

import (
	"prom-lens-backend/internal/apps/authn/controller"
	"prom-lens-backend/internal/apps/authn/repository"
	"prom-lens-backend/internal/apps/authn/service"

	"gorm.io/gorm"
)

type UserProvider struct {
	Repo       repository.UserRepository
	Service    service.UserManager
	Controller *controller.UserController
}

func NewUserProvider(db *gorm.DB) *UserProvider {
	repo := repository.NewUserRepository(db)
	manager := service.NewUserManager(db, repo)
	return &UserProvider{
		Repo:       repo,
		Service:    manager,
		Controller: controller.NewUserController(manager),
	}
}
