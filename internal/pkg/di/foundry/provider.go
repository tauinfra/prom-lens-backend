package foundry

import (
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

type FoundryProvider struct {
	CredGitlab  *CredGitlabProvider
	CredHarbor  *CredHarborProvider
	CredArgoCD  *CredArgoCDProvider
	Project     *ProjectProvider
	Environment *EnvironmentProvider
	Application *ApplicationProvider
	Pipeline    *PipelineProvider
	Release     *ReleaseProvider
	Gitlab      *GitlabProvider
	// ... 其他资源
}

func NewFoundryProvider(db *gorm.DB, encryptor encryption.Encryptor) *FoundryProvider {
	gitlabProvider := NewGitlabProvider(db, encryptor)
	pipelineProvider := NewPipelineProvider(db)
	return &FoundryProvider{
		CredGitlab:  NewCredGitlabProvider(db, encryptor),
		CredHarbor:  NewCredHarborProvider(db, encryptor),
		CredArgoCD:  NewCredArgoCDProvider(db, encryptor),
		Project:     NewProjectProvider(db),
		Environment: NewEnvironmentProvider(db),
		Application: NewApplicationProvider(db),
		Pipeline:    NewPipelineProvider(db),
		Release:     NewReleaseProvider(db, *pipelineProvider.Repo, *gitlabProvider.Service),
		Gitlab:      NewGitlabProvider(db, encryptor),
	}
}
