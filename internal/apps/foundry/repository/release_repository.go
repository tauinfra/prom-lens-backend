package repository

import (
	"context"
	"valyria-backend/internal/apps/foundry/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// ReleaseRepository 定义了数据访问层的接口
type ReleaseRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Release, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Release, error)
	Create(ctx context.Context, data *model.Release) error
	Update(ctx context.Context, id int, data *model.Release) error
	Delete(ctx context.Context, id int) error
	WithTx(tx *gorm.DB) ReleaseRepository
}

// releaseRepository 实现了 ReleaseRepository 接口
type releaseRepository struct {
	db *gorm.DB
}

// NewReleaseRepository 创建新的 ReleaseRepository 实例
func NewReleaseRepository(db *gorm.DB) ReleaseRepository {
	return &releaseRepository{db: db}
}

// List 查询列表
func (r *releaseRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Release, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *releaseRepository) Get(ctx context.Context, id int) (model.Release, error) {
	var data model.Release
	if err := r.db.WithContext(ctx).
		Preload("Pipeline.Environment.Project").
		Preload("Pipeline.Application").
		Preload("Pipeline.CredGitlab").
		First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *releaseRepository) Create(ctx context.Context, data *model.Release) error {
	// 创建
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新
func (r *releaseRepository) Update(ctx context.Context, id int, data *model.Release) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除
func (r *releaseRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Release{}, id).Error
}

// WithTx 返回一个绑定事务的 Repository
func (r *releaseRepository) WithTx(db *gorm.DB) ReleaseRepository {
	return &releaseRepository{db: db}
}
