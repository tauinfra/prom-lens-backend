package service

import (
	"context"
	"valyria-backend/internal/apps/foundry/model"
	"valyria-backend/internal/apps/foundry/repository"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

type GitlabManager interface {
	ListGroups(ctx context.Context, gitlabID int) ([]model.GitlabGroup, error)
	ListGroupProjects(ctx context.Context, gitlabID, groupID int) ([]model.GitlabProject, error)
	ListProjectBranches(ctx context.Context, gitlabID, projectID int) ([]model.GitlabRef, error)
	ListProjectTags(ctx context.Context, gitlabID, projectID int) ([]model.GitlabRef, error)
	GetRef(ctx context.Context, gitlabID, projectID int, ref string) (*model.GitlabRef, error)
}

// gitlabManager 实现 GitlabManager 接口
type gitlabManager struct {
	cred      repository.CredGitlabRepository
	repo      repository.GitlabRepository
	encryptor encryption.Encryptor
	db        *gorm.DB
}

// NewGitlabManager 创建新的 GitlabManager 实例
func NewGitlabManager(cred repository.CredGitlabRepository, repo repository.GitlabRepository, encryptor encryption.Encryptor, db *gorm.DB) GitlabManager {
	return &gitlabManager{cred: cred, repo: repo, encryptor: encryptor, db: db}
}

// ListGroups 列表
func (s *gitlabManager) ListGroups(ctx context.Context, gitlabID int) ([]model.GitlabGroup, error) {
	// 获取凭证
	cred, err := s.cred.Get(ctx, gitlabID)
	if err != nil {
		return nil, err
	}
	// 解密 token
	decryptToken, err := s.encryptor.Decrypt(cred.Token)
	if err != nil {
		return nil, err
	}
	// 查询 Gitlab Groups
	return s.repo.ListGroups(decryptToken, cred.BaseURL)
}

// ListGroupProjects 列表
func (s *gitlabManager) ListGroupProjects(ctx context.Context, gitlabID, groupID int) ([]model.GitlabProject, error) {
	// 获取凭证
	cred, err := s.cred.Get(ctx, gitlabID)
	if err != nil {
		return nil, err
	}
	// 解密 token
	decryptToken, err := s.encryptor.Decrypt(cred.Token)
	if err != nil {
		return nil, err
	}
	// 查询 Gitlab Groups
	return s.repo.ListGroupProjects(decryptToken, cred.BaseURL, groupID)
}

// ListProjectBranches 列表
func (s *gitlabManager) ListProjectBranches(ctx context.Context, gitlabID, projectID int) ([]model.GitlabRef, error) {
	// 获取凭证
	cred, err := s.cred.Get(ctx, gitlabID)
	if err != nil {
		return nil, err
	}
	// 解密 token
	decryptToken, err := s.encryptor.Decrypt(cred.Token)
	if err != nil {
		return nil, err
	}
	// 查询 Gitlab API
	return s.repo.ListProjectBranches(decryptToken, cred.BaseURL, projectID)
}

// ListProjectTags 列表
func (s *gitlabManager) ListProjectTags(ctx context.Context, gitlabID, projectID int) ([]model.GitlabRef, error) {
	// 获取凭证
	cred, err := s.cred.Get(ctx, gitlabID)
	if err != nil {
		return nil, err
	}
	// 解密 token
	decryptToken, err := s.encryptor.Decrypt(cred.Token)
	if err != nil {
		return nil, err
	}
	// 查询 Gitlab API
	return s.repo.ListProjectTags(decryptToken, cred.BaseURL, projectID)
}

func (s *gitlabManager) GetRef(ctx context.Context, gitlabID, projectID int, ref string) (*model.GitlabRef, error) {
	// 获取凭证
	cred, err := s.cred.Get(ctx, gitlabID)
	if err != nil {
		return nil, err
	}
	// 解密 token
	decryptToken, err := s.encryptor.Decrypt(cred.Token)
	if err != nil {
		return nil, err
	}
	// 查询 Gitlab API
	return s.repo.GetRef(decryptToken, cred.BaseURL, ref, projectID)
}
