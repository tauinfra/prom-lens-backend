package repository

import (
	"context"
	"valyria-backend/internal/apps/dragon/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// ReviewRepository 定义了数据访问层的接口
type ReviewRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Review, pg.Pagination, error)
	Get(ctx context.Context, id uint) (model.Review, error)
	Create(ctx context.Context, data *model.Review) error
	Update(ctx context.Context, id uint, data *model.Review) error
	Delete(ctx context.Context, id uint) error
	ExistsReviewerInRelease(ctx context.Context, releaseID uint, reviewerID uint, excludeID uint) (bool, error)
	ExistsStepInRelease(ctx context.Context, releaseID uint, step int) (bool, error)
	WithTx(tx *gorm.DB) ReviewRepository
}

// reviewRepository 实现了 ReviewRepository 接口
type reviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository 创建新的 ReviewRepository 实例
func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

// List 查询列表
func (r *reviewRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Review, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *reviewRepository) Get(ctx context.Context, id uint) (model.Review, error) {
	var data model.Review
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *reviewRepository) Create(ctx context.Context, data *model.Review) error {
	// 创建
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新
func (r *reviewRepository) Update(ctx context.Context, id uint, data *model.Review) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除
func (r *reviewRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Review{}, id).Error
}

func (r *reviewRepository) ExistsReviewerInRelease(ctx context.Context, releaseID uint, reviewerID uint, excludeID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Review{}).
		Where("release_id = ? AND reviewer_id = ? AND id <> ?", releaseID, reviewerID, excludeID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *reviewRepository) ExistsStepInRelease(ctx context.Context, releaseID uint, step int) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Review{}).
		Where("release_id = ? AND step = ?", releaseID, step).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// WithTx 返回一个绑定事务的 Repository
func (r *reviewRepository) WithTx(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}
