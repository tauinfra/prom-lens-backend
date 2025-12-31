package repository

import (
	"context"
	"valyria-backend/internal/apps/foundry/model"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

// CredArgoCDRepository 定义了数据访问层的接口
type CredArgoCDRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.CredArgoCD, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.CredArgoCD, error)
	Create(ctx context.Context, data *model.CredArgoCD) error
	Update(ctx context.Context, id int, data *model.CredArgoCD) error
	Delete(ctx context.Context, id int) error
	WithTx(tx *gorm.DB) CredArgoCDRepository
}

// credArgoCDRepository 实现了 CredArgoCDRepository 接口
type credArgoCDRepository struct {
	db        *gorm.DB
	encryptor encryption.Encryptor
}

// NewCredArgoCDRepository 创建新的 CredArgoCDRepository 实例
func NewCredArgoCDRepository(db *gorm.DB, encryptor encryption.Encryptor) CredArgoCDRepository {
	return &credArgoCDRepository{db: db, encryptor: encryptor}
}

// List 查询列表
func (r *credArgoCDRepository) List(ctx context.Context, params pg.QueryParams) (data []model.CredArgoCD, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *credArgoCDRepository) Get(ctx context.Context, id int) (model.CredArgoCD, error) {
	var data model.CredArgoCD
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *credArgoCDRepository) Create(ctx context.Context, data *model.CredArgoCD) error {
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新
func (r *credArgoCDRepository) Update(ctx context.Context, id int, data *model.CredArgoCD) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除
func (r *credArgoCDRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.CredArgoCD{}, id).Error
}

// WithTx 返回一个绑定事务的 Repository
func (r *credArgoCDRepository) WithTx(db *gorm.DB) CredArgoCDRepository {
	return &credArgoCDRepository{db: db}
}
