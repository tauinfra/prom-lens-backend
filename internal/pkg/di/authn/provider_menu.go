package authn

import (
	"valyria-backend/internal/apps/authn/controller"
	"valyria-backend/internal/apps/authn/repository"
	"valyria-backend/internal/apps/authn/service"

	"gorm.io/gorm"
)

type MenuProvider struct {
	Repo       *repository.MenuRepository
	Service    *service.MenuManager
	Controller *controller.MenuController
}

func NewMenuProvider(db *gorm.DB) *MenuProvider {
	repo := repository.NewMenuRepository(db)
	manager := service.NewMenuManager(db, repo)
	control := controller.NewMenuController(manager)

	return &MenuProvider{
		Repo:       &repo,
		Service:    &manager,
		Controller: control,
	}
}
