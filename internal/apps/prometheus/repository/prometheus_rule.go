package repository

import (
	"context"
	"valyria-backend/internal/apps/prometheus/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// RuleRepository 定义接口
type RuleRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Rule, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Rule, error)
	Create(ctx context.Context, data *model.Rule) error
	Update(ctx context.Context, id int, data *model.Rule) error
	Delete(ctx context.Context, id int) error
	WithTx(tx *gorm.DB) RuleRepository
}

// ruleRepository 实现了 RuleRepository 接口
type ruleRepository struct {
	db *gorm.DB
}

// NewRuleRepository 创建新的 RuleRepository 实例
func NewRuleRepository(db *gorm.DB) RuleRepository {
	return &ruleRepository{db: db}
}

// List 列表
func (r *ruleRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Rule, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *ruleRepository) Get(ctx context.Context, id int) (data model.Rule, err error) {
	if err = r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return
}

// Create 创建
func (r *ruleRepository) Create(ctx context.Context, data *model.Rule) error {
	return r.db.WithContext(ctx).Create(data).Error
}

// Update 更新
func (r *ruleRepository) Update(ctx context.Context, id int, data *model.Rule) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error
}

// Delete 删除
func (r *ruleRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Rule{}).Error
}

func (r *ruleRepository) WithTx(db *gorm.DB) RuleRepository {
	return &ruleRepository{db: db}
}
