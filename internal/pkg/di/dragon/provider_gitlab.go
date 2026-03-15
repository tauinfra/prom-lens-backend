package dragon

import (
	"valyria-backend/internal/apps/dragon/controller"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/service"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

type GitlabProvider struct {
	Repo       *repository.GitlabRepository
	Service    *service.GitlabManager
	Controller *controller.GitlabController
}

func NewGitlabProvider(db *gorm.DB, encryptor encryption.Encryptor) *GitlabProvider {
	credentialRepo := repository.NewCredentialRepository(db, encryptor)
	repo := repository.NewGitlabRepository()
	svc := service.NewGitlabManager(credentialRepo, repo, encryptor, db)
	ctrl := controller.NewGitlabController(svc)

	return &GitlabProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
