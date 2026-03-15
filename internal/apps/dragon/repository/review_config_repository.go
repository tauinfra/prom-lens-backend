package repository

import (
	"context"
	"errors"
	"valyria-backend/internal/apps/dragon/dto"
	"valyria-backend/internal/apps/dragon/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// ReviewConfigRepository 审批配置数据访问层
type ReviewConfigRepository interface {
	ListDTO(ctx context.Context, params pg.QueryParams) ([]dto.ReviewStageDTO, pg.Pagination, error)
	GetDTO(ctx context.Context, id uint) (dto.ReviewStageDTO, error)
	Create(ctx context.Context, data *model.ReviewStageConfig) error
	Update(ctx context.Context, id uint, data *model.ReviewStageConfig) error
	Delete(ctx context.Context, id uint) error
	ListAll(ctx context.Context) ([]model.ReviewStageConfig, error)
	UserHasRole(ctx context.Context, userID, roleID uint) (bool, error)
	WithTx(tx *gorm.DB) ReviewConfigRepository
}

type reviewConfigRepository struct {
	db *gorm.DB
}

func NewReviewConfigRepository(db *gorm.DB) ReviewConfigRepository {
	return &reviewConfigRepository{db: db}
}

func (r *reviewConfigRepository) ListDTO(ctx context.Context, params pg.QueryParams) (data []dto.ReviewStageDTO, pagination pg.Pagination, err error) {
	query := r.db.WithContext(ctx).
		Table("valyria_dragon_review_stage s").
		Select(`
        s.id,
        s.step,
        s.role_id,
        r.name AS role_name,
        s.creator,
        s.created_at,
        s.updated_at
    `).
		Joins("LEFT JOIN valyria_authn_role r ON r.id = s.role_id")
	if params.SortBy == "" {
		params.SortBy = "step"
	}
	if pagination, err = pg.Paginate(query, &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

func (r *reviewConfigRepository) GetDTO(ctx context.Context, id uint) (dto.ReviewStageDTO, error) {
	var data dto.ReviewStageDTO
	if err := r.db.WithContext(ctx).
		Table("valyria_dragon_review_stage s").
		Select(`
        s.id,
        s.step,
        s.role_id,
        r.name AS role_name,
        s.creator,
        s.created_at,
        s.updated_at
    `).
		Joins("LEFT JOIN valyria_authn_role r ON r.id = s.role_id").
		Where("s.id = ?", id).
		First(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (r *reviewConfigRepository) Create(ctx context.Context, data *model.ReviewStageConfig) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *reviewConfigRepository) Update(ctx context.Context, id uint, data *model.ReviewStageConfig) error {
	return r.db.WithContext(ctx).Model(&model.ReviewStageConfig{}).Where("id = ?", id).Updates(data).Error
}

func (r *reviewConfigRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.ReviewStageConfig{}, id).Error
}

func (r *reviewConfigRepository) ListAll(ctx context.Context) ([]model.ReviewStageConfig, error) {
	var data []model.ReviewStageConfig
	err := r.db.WithContext(ctx).
		Preload("Role").
		Order("step asc").
		Find(&data).Error
	return data, err
}

func (r *reviewConfigRepository) UserHasRole(ctx context.Context, userID, roleID uint) (bool, error) {
	if userID == 0 || roleID == 0 {
		return false, errors.New("invalid role check parameters")
	}
	var count int64
	err := r.db.WithContext(ctx).
		Table("valyria_authn_user_role ur").
		Where("ur.user_id = ? AND ur.role_id = ?", userID, roleID).
		Count(&count).Error
	return count > 0, err
}

func (r *reviewConfigRepository) WithTx(db *gorm.DB) ReviewConfigRepository {
	return &reviewConfigRepository{db: db}
}
