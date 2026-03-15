package dragon

import (
	"valyria-backend/internal/apps/dragon/controller"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/service"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

type CredentialProvider struct {
	Repo       *repository.CredentialRepository
	Service    *service.CredentialManager
	Controller *controller.CredentialController
}

func NewCredentialProvider(db *gorm.DB, encryptor encryption.Encryptor) *CredentialProvider {
	repo := repository.NewCredentialRepository(db, encryptor)
	svc := service.NewCredentialManager(repo, encryptor, db)
	control := controller.NewCredentialController(svc)

	return &CredentialProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: control,
	}
}
