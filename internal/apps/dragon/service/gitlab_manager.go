package service

import (
	"context"
	"valyria-backend/internal/apps/dragon/model"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

type GitlabManager interface {
	ListGroups(ctx context.Context, gitlabID int64) ([]model.GitlabGroup, error)
	ListGroupProjects(ctx context.Context, gitlabID, groupID int64) ([]model.GitlabProject, error)
	ListProjectBranches(ctx context.Context, gitlabID, projectID int64) ([]model.GitlabRef, error)
	ListProjectTags(ctx context.Context, gitlabID, projectID int64) ([]model.GitlabRef, error)
	GetRef(ctx context.Context, gitlabID, projectID int64, ref string) (*model.GitlabRef, error)
	GetRepoURL(ctx context.Context, gitlabID, projectID int64) (string, error)
}

// gitlabManager 实现 GitlabManager 接口
type gitlabManager struct {
	credentialRepo repository.CredentialRepository
	repo           repository.GitlabRepository
	encryptor      encryption.Encryptor
	db             *gorm.DB
}

// NewGitlabManager 创建新的 GitlabManager 实例
func NewGitlabManager(
	credentialRepo repository.CredentialRepository,
	repo repository.GitlabRepository,
	encryptor encryption.Encryptor,
	db *gorm.DB) GitlabManager {
	return &gitlabManager{credentialRepo: credentialRepo, repo: repo, encryptor: encryptor, db: db}
}

// getGitlabCredentials 获取并解密 GitLab 凭证
func (s *gitlabManager) getGitlabCredentials(ctx context.Context, gitlabID int64) (token, baseURL string, err error) {
	// 获取凭证
	credential, err := s.credentialRepo.Get(ctx, uint(gitlabID))
	if err != nil {
		return "", "", err
	}
	// 解密 token
	decryptedToken, err := s.encryptor.Decrypt(credential.EncryptedSecret)
	if err != nil {
		return "", "", err
	}
	return decryptedToken, credential.BaseURL, nil
}

// ListGroups 列表
func (s *gitlabManager) ListGroups(ctx context.Context, gitlabID int64) ([]model.GitlabGroup, error) {
	token, baseURL, err := s.getGitlabCredentials(ctx, gitlabID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListGroups(token, baseURL)
}

// ListGroupProjects 列表
func (s *gitlabManager) ListGroupProjects(ctx context.Context, gitlabID, groupID int64) ([]model.GitlabProject, error) {
	token, baseURL, err := s.getGitlabCredentials(ctx, gitlabID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListGroupProjects(token, baseURL, groupID)
}

// ListProjectBranches 列表
func (s *gitlabManager) ListProjectBranches(ctx context.Context, gitlabID, projectID int64) ([]model.GitlabRef, error) {
	token, baseURL, err := s.getGitlabCredentials(ctx, gitlabID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListProjectBranches(token, baseURL, projectID)
}

// ListProjectTags 列表
func (s *gitlabManager) ListProjectTags(ctx context.Context, gitlabID, projectID int64) ([]model.GitlabRef, error) {
	token, baseURL, err := s.getGitlabCredentials(ctx, gitlabID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListProjectTags(token, baseURL, projectID)
}

func (s *gitlabManager) GetRef(ctx context.Context, gitlabID, projectID int64, ref string) (*model.GitlabRef, error) {
	token, baseURL, err := s.getGitlabCredentials(ctx, gitlabID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetRef(token, baseURL, ref, projectID)
}

func (s *gitlabManager) GetRepoURL(ctx context.Context, gitlabID, projectID int64) (string, error) {
	token, baseURL, err := s.getGitlabCredentials(ctx, gitlabID)
	if err != nil {
		return "", err
	}
	return s.repo.GetRepoURL(token, baseURL, projectID)
}
