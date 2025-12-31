package foundry

import (
	"valyria-backend/internal/apps/foundry/controller"
	"valyria-backend/internal/apps/foundry/repository"
	"valyria-backend/internal/apps/foundry/service"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

type GitlabProvider struct {
	Repo       *repository.GitlabRepository
	Service    *service.GitlabManager
	Controller *controller.GitlabController
}

func NewGitlabProvider(db *gorm.DB, encryptor encryption.Encryptor) *GitlabProvider {
	cred := repository.NewCredGitlabRepository(db, encryptor)
	repo := repository.NewGitlabRepository()
	svc := service.NewGitlabManager(cred, repo, encryptor, db)
	ctrl := controller.NewGitlabController(svc)

	return &GitlabProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
