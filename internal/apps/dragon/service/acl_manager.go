package service

import (
	"context"
	"errors"
	"valyria-backend/internal/apps/dragon/dto"
	"valyria-backend/internal/apps/dragon/model"
	"valyria-backend/internal/apps/dragon/repository"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type PipelineACLManager interface {
	List(ctx context.Context, pipelineID uint, params pg.QueryParams) ([]dto.PipelineACLDTO, pg.Pagination, error)
	Create(ctx context.Context, pipelineID uint, userID uint, action string, creator string) error
	BatchCreate(ctx context.Context, pipelineID uint, users []uint, action string, creator string) error
	Delete(ctx context.Context, id uint) error
	HasAction(ctx context.Context, userID, pipelineID uint, action string) (bool, error)
	ListPipelineIDs(ctx context.Context, userID uint, action string) ([]uint, error)
}

type pipelineACLManager struct {
	repo repository.PipelineACLRepository
}

func NewPipelineACLManager(repo repository.PipelineACLRepository) PipelineACLManager {
	return &pipelineACLManager{repo: repo}
}

func (s *pipelineACLManager) List(ctx context.Context, pipelineID uint, params pg.QueryParams) ([]dto.PipelineACLDTO, pg.Pagination, error) {
	return s.repo.ListDTOByPipeline(ctx, pipelineID, params)
}

func (s *pipelineACLManager) Create(ctx context.Context, pipelineID uint, userID uint, action string, creator string) error {
	if pipelineID == 0 || userID == 0 || action == "" {
		return gorm.ErrInvalidData
	}
	if !IsValidAction(action) {
		return errors.New("invalid action")
	}
	data := &model.PipelineACL{
		PipelineID: pipelineID,
		UserID:     userID,
		Action:     action,
		Creator:    creator,
	}
	return s.repo.Create(ctx, data)
}

func (s *pipelineACLManager) BatchCreate(ctx context.Context, pipelineID uint, users []uint, action string, creator string) error {
	if pipelineID == 0 || len(users) == 0 || action == "" {
		return gorm.ErrInvalidData
	}
	if !IsValidAction(action) {
		return errors.New("invalid action")
	}
	unique := make(map[uint]struct{}, len(users))
	data := make([]model.PipelineACL, 0, len(users))
	for _, userID := range users {
		if userID == 0 {
			continue
		}
		if _, exists := unique[userID]; exists {
			continue
		}
		unique[userID] = struct{}{}
		data = append(data, model.PipelineACL{
			PipelineID: pipelineID,
			UserID:     userID,
			Action:     action,
			Creator:    creator,
		})
	}
	if len(data) == 0 {
		return gorm.ErrInvalidData
	}
	return s.repo.BatchCreate(ctx, data)
}

func (s *pipelineACLManager) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return gorm.ErrInvalidData
	}
	return s.repo.Delete(ctx, id)
}

func (s *pipelineACLManager) HasAction(ctx context.Context, userID, pipelineID uint, action string) (bool, error) {
	if pipelineID == 0 || userID == 0 || action == "" {
		return false, errors.New("invalid acl parameters")
	}
	if IsSuperuser(ctx) {
		return true, nil
	}
	if action == ActionView {
		return s.hasAnyAction(ctx, userID, pipelineID, []string{
			ActionView, ActionDeploy, ActionApprove, ActionRollback,
		})
	}
	return s.repo.HasUserAction(ctx, userID, pipelineID, action)
}

func (s *pipelineACLManager) ListPipelineIDs(ctx context.Context, userID uint, action string) ([]uint, error) {
	if userID == 0 || action == "" {
		return nil, errors.New("invalid acl parameters")
	}
	return s.repo.ListPipelineIDs(ctx, userID, action)
}

func (s *pipelineACLManager) hasAnyAction(ctx context.Context, userID, pipelineID uint, actions []string) (bool, error) {
	for _, action := range actions {
		ok, err := s.repo.HasUserAction(ctx, userID, pipelineID, action)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// isSuperuser moved to context_utils.go
