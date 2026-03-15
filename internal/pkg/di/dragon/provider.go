package dragon

import (
	"valyria-backend/internal/core/config"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

type Provider struct {
	Project     *ProjectProvider
	Environment *EnvironmentProvider
	Pipeline    *PipelineProvider
	TektonPipeline *TektonPipelineProvider
	TektonTask    *TektonTaskProvider
	TektonTaskRun *TektonTaskRunProvider
	TektonPipelineRun *TektonPipelineRunProvider
	Release     *ReleaseProvider
	Report      *ReportProvider
	ReviewConfig *ReviewConfigProvider
	Review      *ReviewProvider
	PipelineACL *PipelineACLProvider
	Credential  *CredentialProvider
	Gitlab      *GitlabProvider
	// ... 其他资源
}

func NewDragonProvider(db *gorm.DB, encryptor encryption.Encryptor, cfg *config.Config) *Provider {
	credentialProvider := NewCredentialProvider(db, encryptor)
	gitlabProvider := NewGitlabProvider(db, encryptor)
	aclProvider := NewPipelineACLProvider(db)
	projectProvider := NewProjectProvider(db, *aclProvider.Repo)
	environmentProvider := NewEnvironmentProvider(db, *aclProvider.Repo, *credentialProvider.Repo)
	pipelineProvider := NewPipelineProvider(db, *credentialProvider.Repo, *aclProvider.Repo)
	tektonPipelineProvider := NewTektonPipelineProvider(cfg)
	tektonTaskProvider := NewTektonTaskProvider(cfg)
	tektonTaskRunProvider := NewTektonTaskRunProvider(cfg)
	tektonPipelineRunProvider := NewTektonPipelineRunProvider(cfg)
	reviewConfigProvider := NewReviewConfigProvider(db)
	reviewProvider := NewReviewProvider(db, *reviewConfigProvider.Repo, *aclProvider.Repo, *gitlabProvider.Service, cfg)

	return &Provider{
		Project:      projectProvider,
		Environment:  environmentProvider,
		Pipeline:     pipelineProvider,
		TektonPipeline: tektonPipelineProvider,
		TektonTask: tektonTaskProvider,
		TektonTaskRun: tektonTaskRunProvider,
		TektonPipelineRun: tektonPipelineRunProvider,
		Release:      NewReleaseProvider(db, *pipelineProvider.Repo, *aclProvider.Repo, *reviewProvider.Service, *gitlabProvider.Service, cfg),
		Report:       NewReportProvider(db),
		ReviewConfig: reviewConfigProvider,
		Review:       reviewProvider,
		PipelineACL:  aclProvider,
		Credential:   credentialProvider,
		Gitlab:       gitlabProvider,
	}
}
