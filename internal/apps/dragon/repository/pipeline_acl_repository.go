package repository

import (
	"context"
	"valyria-backend/internal/apps/dragon/dto"
	"valyria-backend/internal/apps/dragon/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm/clause"
	"gorm.io/gorm"
)

type PipelineACLRepository interface {
	ListDTOByPipeline(ctx context.Context, pipelineID uint, params pg.QueryParams) ([]dto.PipelineACLDTO, pg.Pagination, error)
	Create(ctx context.Context, data *model.PipelineACL) error
	BatchCreate(ctx context.Context, data []model.PipelineACL) error
	Delete(ctx context.Context, id uint) error
	HasUserAction(ctx context.Context, userID, pipelineID uint, action string) (bool, error)
	ListPipelineIDs(ctx context.Context, userID uint, action string) ([]uint, error)
	ListProjectIDs(ctx context.Context, userID uint, action string) ([]uint, error)
	ListEnvironmentIDs(ctx context.Context, userID uint, projectID uint, action string) ([]uint, error)
	WithTx(tx *gorm.DB) PipelineACLRepository
}

type pipelineACLRepository struct {
	db *gorm.DB
}

func NewPipelineACLRepository(db *gorm.DB) PipelineACLRepository {
	return &pipelineACLRepository{db: db}
}

func (r *pipelineACLRepository) ListDTOByPipeline(ctx context.Context, pipelineID uint, params pg.QueryParams) (data []dto.PipelineACLDTO, pagination pg.Pagination, err error) {
	query := r.db.WithContext(ctx).
		Table("valyria_dragon_pipeline_acl a").
		Select(`
			a.id,
			a.pipeline_id,
			p.name AS pipeline_name,
			a.user_id,
			u.username AS user_name,
			a.action,
			a.creator,
			a.created_at,
			a.updated_at
		`).
		Joins("JOIN valyria_dragon_pipeline p ON p.id = a.pipeline_id").
		Joins("JOIN valyria_authn_user u ON u.id = a.user_id").
		Where("a.pipeline_id = ?", pipelineID)
	if pagination, err = pg.Paginate(query, &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

func (r *pipelineACLRepository) Create(ctx context.Context, data *model.PipelineACL) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *pipelineACLRepository) BatchCreate(ctx context.Context, data []model.PipelineACL) error {
	if len(data) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&data).Error
}

func (r *pipelineACLRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.PipelineACL{}, id).Error
}

func (r *pipelineACLRepository) HasUserAction(ctx context.Context, userID, pipelineID uint, action string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.PipelineACL{}).
		Where("user_id = ? AND pipeline_id = ? AND action = ?", userID, pipelineID, action).
		Count(&count).Error
	return count > 0, err
}

func (r *pipelineACLRepository) ListPipelineIDs(ctx context.Context, userID uint, action string) ([]uint, error) {
	var ids []uint
	err := r.db.WithContext(ctx).
		Model(&model.PipelineACL{}).
		Where("user_id = ? AND action = ?", userID, action).
		Pluck("pipeline_id", &ids).Error
	return ids, err
}

func (r *pipelineACLRepository) ListProjectIDs(ctx context.Context, userID uint, action string) ([]uint, error) {
	var ids []uint
	err := r.db.WithContext(ctx).
		Table("valyria_dragon_pipeline_acl a").
		Select("DISTINCT p.id").
		Joins("JOIN valyria_dragon_pipeline pl ON pl.id = a.pipeline_id").
		Joins("JOIN valyria_dragon_environment e ON e.id = pl.environment_id").
		Joins("JOIN valyria_dragon_project p ON p.id = e.project_id").
		Where("a.user_id = ? AND a.action = ?", userID, action).
		Pluck("p.id", &ids).Error
	return ids, err
}

func (r *pipelineACLRepository) ListEnvironmentIDs(ctx context.Context, userID uint, projectID uint, action string) ([]uint, error) {
	var ids []uint
	query := r.db.WithContext(ctx).
		Table("valyria_dragon_pipeline_acl a").
		Select("DISTINCT e.id").
		Joins("JOIN valyria_dragon_pipeline pl ON pl.id = a.pipeline_id").
		Joins("JOIN valyria_dragon_environment e ON e.id = pl.environment_id").
		Where("a.user_id = ? AND a.action = ?", userID, action)
	if projectID > 0 {
		query = query.Where("e.project_id = ?", projectID)
	}
	err := query.Pluck("e.id", &ids).Error
	return ids, err
}

func (r *pipelineACLRepository) WithTx(db *gorm.DB) PipelineACLRepository {
	return &pipelineACLRepository{db: db}
}
