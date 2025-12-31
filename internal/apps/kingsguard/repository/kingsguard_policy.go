package repository

import (
	"context"
	"valyria-backend/internal/apps/kingsguard/model"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/kingsguard"

	"gorm.io/gorm"
)

// PolicyRepository 定义接口
type PolicyRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.CasbinPolicy, pg.Pagination, error)
	Create(data model.CasbinPolicy) (bool, error)
	Delete(ctx context.Context, id int) (bool, error)
	WithTx(tx *gorm.DB) PolicyRepository
}

// policyRepository 实现了 PolicyRepository 接口
type policyRepository struct {
	db *gorm.DB
}

// NewPolicyRepository 创建新的 PolicyRepository 实例
func NewPolicyRepository(db *gorm.DB) PolicyRepository {
	return &policyRepository{db: db}
}

// List 列表
func (r *policyRepository) List(ctx context.Context, params pg.QueryParams) (data []model.CasbinPolicy, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

func (r *policyRepository) Create(data model.CasbinPolicy) (bool, error) {
	return kingsguard.Enforcer.AddPolicy(data.V0, data.V1, data.V2, data.V3, data.V4)
}

func (r *policyRepository) Delete(ctx context.Context, id int) (bool, error) {
	var data model.CasbinPolicy
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return false, err
	}
	return kingsguard.Enforcer.RemovePolicy(data.V0, data.V1, data.V2, data.V3, data.V4)
}

// WithTx 返回一个绑定事务的 Repository
func (r *policyRepository) WithTx(db *gorm.DB) PolicyRepository {
	return &policyRepository{db: db}
}
