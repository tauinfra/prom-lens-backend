package foundry

import (
	"valyria-backend/internal/apps/foundry/controller"
	"valyria-backend/internal/apps/foundry/repository"
	"valyria-backend/internal/apps/foundry/service"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

type CredGitlabProvider struct {
	Repo       *repository.CredGitlabRepository
	Service    *service.CredGitlabManager
	Controller *controller.CredGitlabController
}

type CredHarborProvider struct {
	Repo       *repository.CredHarborRepository
	Service    *service.CredHarborManager
	Controller *controller.CredHarborController
}

type CredArgoCDProvider struct {
	Repo       *repository.CredArgoCDRepository
	Service    *service.CredArgoCDManager
	Controller *controller.CredArgoCDController
}

func NewCredGitlabProvider(db *gorm.DB, encryptor encryption.Encryptor) *CredGitlabProvider {
	repo := repository.NewCredGitlabRepository(db, encryptor)
	svc := service.NewCredGitlabManager(repo, db)
	ctrl := controller.NewCredGitlabController(svc)

	return &CredGitlabProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}

func NewCredHarborProvider(db *gorm.DB, encryptor encryption.Encryptor) *CredHarborProvider {
	repo := repository.NewCredHarborRepository(db, encryptor)
	svc := service.NewCredHarborManager(repo, db)
	ctrl := controller.NewCredHarborController(svc)

	return &CredHarborProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}

func NewCredArgoCDProvider(db *gorm.DB, encryptor encryption.Encryptor) *CredArgoCDProvider {
	repo := repository.NewCredArgoCDRepository(db, encryptor)
	svc := service.NewCredArgoCDManager(repo, db)
	ctrl := controller.NewCredArgoCDController(svc)

	return &CredArgoCDProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
