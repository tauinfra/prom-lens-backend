package service

import (
	"context"
	"valyria-backend/internal/apps/dragon/dto"
	"valyria-backend/internal/apps/dragon/model"
	"valyria-backend/internal/apps/dragon/repository"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

type CredentialManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.CredentialDTO, pg.Pagination, error)
	Get(ctx context.Context, id uint) (dto.CredentialDTO, error)
	Create(ctx context.Context, data *model.Credential) error
	Update(ctx context.Context, id uint, data *model.Credential) error
	Delete(ctx context.Context, id uint) error
}

// CredentialManager 实现 CredentialManager 接口
type credentialManager struct {
	repo      repository.CredentialRepository
	encryptor encryption.Encryptor
	db        *gorm.DB
}

// NewCredentialManager 创建新的 CredentialManager 实例
func NewCredentialManager(repo repository.CredentialRepository, encryptor encryption.Encryptor, db *gorm.DB) CredentialManager {
	return &credentialManager{repo: repo, encryptor: encryptor, db: db}
}

// List 列表
func (s *credentialManager) List(ctx context.Context, params pg.QueryParams) ([]dto.CredentialDTO, pg.Pagination, error) {
	data, pagination, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, pg.Pagination{}, err
	}
	result := make([]dto.CredentialDTO, 0, len(data))
	for _, item := range data {
		result = append(result, dto.ToCredentialDTO(item))
	}
	return result, pagination, nil
}

// Get 查询
func (s *credentialManager) Get(ctx context.Context, id uint) (dto.CredentialDTO, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil {
		return dto.CredentialDTO{}, err
	}
	return dto.ToCredentialDTO(item), nil
}

// Create 创建
func (s *credentialManager) Create(ctx context.Context, data *model.Credential) error {
	// 加密敏感字段
	encryptedSecret, err := s.encryptor.Encrypt(data.EncryptedSecret)
	if err != nil {
		return err
	}
	data.EncryptedSecret = encryptedSecret
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *credentialManager) Update(ctx context.Context, id uint, data *model.Credential) error {
	// 1. 从数据库读取原始数据
	credential, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	// 2. 判断是否需要更新 Secret
	if data.EncryptedSecret != "" && data.EncryptedSecret != credential.EncryptedSecret {
		encryptedSecret, err := s.encryptor.Encrypt(data.EncryptedSecret)
		if err != nil {
			return err
		}
		data.EncryptedSecret = encryptedSecret
	}
	return s.repo.Update(ctx, id, data)
}

// Delete 更新
func (s *credentialManager) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

