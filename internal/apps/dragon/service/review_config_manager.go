package service

import (
	"context"
	"valyria-backend/internal/apps/dragon/dto"
	"valyria-backend/internal/apps/dragon/model"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/request"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type ReviewScope struct {
	ProjectID     uint
	EnvironmentID uint
	PipelineID    uint
}

type ReviewConfigManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.ReviewStageDTO, pg.Pagination, error)
	Get(ctx context.Context, id uint) (dto.ReviewStageDTO, error)
	Create(ctx context.Context, req *request.CreateReviewStageRequest, creator string) error
	Update(ctx context.Context, id uint, req *request.UpdateReviewStageRequest) error
	Delete(ctx context.Context, id uint) error
}

type reviewConfigManager struct {
	repo repository.ReviewConfigRepository
	db   *gorm.DB
}

func NewReviewConfigManager(
	repo repository.ReviewConfigRepository,
	db *gorm.DB,
) ReviewConfigManager {
	return &reviewConfigManager{repo: repo, db: db}
}

func (s *reviewConfigManager) List(ctx context.Context, params pg.QueryParams) ([]dto.ReviewStageDTO, pg.Pagination, error) {
	return s.repo.ListDTO(ctx, params)
}

func (s *reviewConfigManager) Get(ctx context.Context, id uint) (dto.ReviewStageDTO, error) {
	return s.repo.GetDTO(ctx, id)
}

func (s *reviewConfigManager) Create(ctx context.Context, req *request.CreateReviewStageRequest, creator string) error {
	if req.Step <= 0 {
		return gorm.ErrInvalidData
	}
	data := &model.ReviewStageConfig{
		Step:       req.Step,
		RoleID:     req.RoleID,
		Creator:    creator,
	}
	return s.repo.Create(ctx, data)
}

func (s *reviewConfigManager) Update(ctx context.Context, id uint, req *request.UpdateReviewStageRequest) error {
	update := &model.ReviewStageConfig{}
	if req.Step != nil {
		if *req.Step <= 0 {
			return gorm.ErrInvalidData
		}
		update.Step = *req.Step
	}
	if req.RoleID != nil {
		update.RoleID = *req.RoleID
	}
	if update.Step == 0 && update.RoleID == 0 {
		return gorm.ErrInvalidData
	}
	return s.repo.Update(ctx, id, update)
}

func (s *reviewConfigManager) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
