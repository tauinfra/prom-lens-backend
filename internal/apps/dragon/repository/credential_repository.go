package repository

import (
	"context"
	"valyria-backend/internal/apps/dragon/model"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

// CredentialRepository 定义了数据访问层的接口
type CredentialRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Credential, pg.Pagination, error)
	Get(ctx context.Context, id uint) (model.Credential, error)
	Create(ctx context.Context, data *model.Credential) error
	Update(ctx context.Context, id uint, data *model.Credential) error
	Delete(ctx context.Context, id uint) error
	WithTx(tx *gorm.DB) CredentialRepository
}

// credentialRepository 实现了 CredentialRepository 接口
type credentialRepository struct {
	db        *gorm.DB
	encryptor encryption.Encryptor
}

// NewCredentialRepository 创建新的 CredentialRepository 实例
func NewCredentialRepository(db *gorm.DB, encryptor encryption.Encryptor) CredentialRepository {
	return &credentialRepository{db: db, encryptor: encryptor}
}

// List 查询列表
func (r *credentialRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Credential, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *credentialRepository) Get(ctx context.Context, id uint) (model.Credential, error) {
	var data model.Credential
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *credentialRepository) Create(ctx context.Context, data *model.Credential) error {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新
func (r *credentialRepository) Update(ctx context.Context, id uint, data *model.Credential) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除
func (r *credentialRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Credential{}, id).Error
}

// WithTx 返回一个绑定事务的 Repository（保留 encryptor 以便事务内加解密）
func (r *credentialRepository) WithTx(db *gorm.DB) CredentialRepository {
	return &credentialRepository{db: db, encryptor: r.encryptor}
}
