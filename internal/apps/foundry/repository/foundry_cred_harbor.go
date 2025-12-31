package repository

import (
	"context"
	"gorm.io/gorm"
	"valyria-backend/internal/apps/foundry/model"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/encryption"
)

// CredHarborRepository 定义了数据访问层的接口
type CredHarborRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.CredHarbor, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.CredHarbor, error)
	Create(ctx context.Context, data *model.CredHarbor) error
	Update(ctx context.Context, id int, data *model.CredHarbor) error
	Delete(ctx context.Context, id int) error
	WithTx(tx *gorm.DB) CredHarborRepository
}

// CredHarborRepository 实现了 CredHarborRepository 接口
type credHarborRepository struct {
	db        *gorm.DB
	encryptor encryption.Encryptor
}

// NewCredHarborRepository 创建新的 CredHarborRepository 实例
func NewCredHarborRepository(db *gorm.DB, encryptor encryption.Encryptor) CredHarborRepository {
	return &credHarborRepository{db: db, encryptor: encryptor}
}

// List 查询列表
func (r *credHarborRepository) List(ctx context.Context, params pg.QueryParams) (data []model.CredHarbor, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *credHarborRepository) Get(ctx context.Context, id int) (model.CredHarbor, error) {
	var data model.CredHarbor
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *credHarborRepository) Create(ctx context.Context, data *model.CredHarbor) error {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新
func (r *credHarborRepository) Update(ctx context.Context, id int, data *model.CredHarbor) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除
func (r *credHarborRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.CredHarbor{}, id).Error
}

// WithTx 返回一个绑定事务的 Repository
func (r *credHarborRepository) WithTx(db *gorm.DB) CredHarborRepository {
	return &credHarborRepository{db: db}
}
