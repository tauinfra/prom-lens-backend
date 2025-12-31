package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"valyria-backend/internal/apps/foundry/model"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/encryption"
)

// CredGitlabRepository 定义了数据访问层的接口
type CredGitlabRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.CredGitlab, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.CredGitlab, error)
	Create(ctx context.Context, data *model.CredGitlab) error
	Update(ctx context.Context, id int, data *model.CredGitlab) error
	Delete(ctx context.Context, id int) error
	WithTx(tx *gorm.DB) CredGitlabRepository
}

// CredGitlabRepository 实现了 CredGitlabRepository 接口
type credGitlabRepository struct {
	db        *gorm.DB
	encryptor encryption.Encryptor
}

// NewCredGitlabRepository 创建新的 CredGitlabRepository 实例
func NewCredGitlabRepository(db *gorm.DB, encryptor encryption.Encryptor) CredGitlabRepository {
	return &credGitlabRepository{db: db, encryptor: encryptor}
}

// List 查询列表
func (r *credGitlabRepository) List(ctx context.Context, params pg.QueryParams) (data []model.CredGitlab, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *credGitlabRepository) Get(ctx context.Context, id int) (model.CredGitlab, error) {
	var data model.CredGitlab
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *credGitlabRepository) Create(ctx context.Context, data *model.CredGitlab) error {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新
func (r *credGitlabRepository) Update(ctx context.Context, id int, data *model.CredGitlab) error {
	// 1. 先查询现有数据
	var cluster model.CredGitlab
	if err := r.db.WithContext(ctx).First(&cluster, id).Error; err != nil {
		return err
	}
	// 2. 检查 Token 是否有变化
	encryptTokenChanged := data.Token != "" && data.Token != cluster.Token
	if encryptTokenChanged {
		encryptToken, err := r.encryptor.Encrypt(data.Token)
		if err != nil {
			return fmt.Errorf("gitlab encrypt token failed, err: %v", err)
		}
		data.Token = encryptToken
	}

	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除
func (r *credGitlabRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.CredGitlab{}, id).Error
}

// WithTx 返回一个绑定事务的 Repository
func (r *credGitlabRepository) WithTx(db *gorm.DB) CredGitlabRepository {
	return &credGitlabRepository{db: db}
}
