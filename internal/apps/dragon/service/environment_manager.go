package service

import (
	"context"
	"errors"
	"valyria-backend/internal/apps/dragon/dto"
	"valyria-backend/internal/apps/dragon/model"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/request"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type EnvironmentManager interface {
	List(ctx context.Context, userID, projectID uint, isAdmin bool, params pg.QueryParams) ([]dto.EnvironmentDTO, pg.Pagination, error)
	Get(ctx context.Context, projectID, environmentID uint) (dto.EnvironmentDTO, error)
	Create(ctx context.Context, projectID uint, req *request.CreateEnvironmentRequest) error
	Update(ctx context.Context, projectID, environmentID uint, req *request.UpdateEnvironmentRequest) error
	Delete(ctx context.Context, projectID, environmentID uint) error
}

// EnvironmentManager 实现 EnvironmentManager 接口
type environmentManager struct {
	repo      repository.EnvironmentRepository
	aclRepo   repository.PipelineACLRepository
	credRepo  repository.CredentialRepository
	db        *gorm.DB
}

// NewEnvironmentManager 创建新的 EnvironmentManager 实例
func NewEnvironmentManager(
	repo repository.EnvironmentRepository,
	aclRepo repository.PipelineACLRepository,
	credRepo repository.CredentialRepository,
	db *gorm.DB,
) EnvironmentManager {
	return &environmentManager{
		repo:      repo,
		aclRepo:   aclRepo,
		credRepo:  credRepo,
		db:        db,
	}
}

// List 列表
func (s *environmentManager) List(ctx context.Context, userID, projectID uint, isAdmin bool, params pg.QueryParams) (data []dto.EnvironmentDTO, pagination pg.Pagination, err error) {
	var ids []uint
	if !isAdmin {
		if userID == 0 {
			return []dto.EnvironmentDTO{}, pagination, errors.New("no permission")
		}
		if ids, err = s.listViewableEnvironmentIDs(ctx, userID, projectID); err != nil {
			return
		}
		if len(ids) == 0 {
			return []dto.EnvironmentDTO{}, pg.Pagination{Page: params.Page, Size: params.Size, Total: 0}, nil
		}
		params.Conditions = []pg.Condition{{Field: "id", Operator: "IN", Value: ids}}
	}

	environments, pagination, err := s.repo.List(ctx, params)
	if err != nil {
		return data, pagination, err
	}
	for _, env := range environments {
		data = append(data, dto.EnvironmentDTO{
			ID:            env.ID,
			Name:          env.Name,
			ProjectID:     env.ProjectID,
			ProjectName:   env.Project.Name,
			HarborID:      env.HarborID,
			HarborBaseURL: env.Harbor.BaseURL,
			Description:   env.Description,
			Creator:       env.Creator,
			CreatedAt:     env.CreatedAt,
			UpdatedAt:     env.UpdatedAt,
		})
	}
	return
}

func (s *environmentManager) listViewableEnvironmentIDs(ctx context.Context, userID, projectID uint) ([]uint, error) {
	actions := []string{ActionView, ActionDeploy, ActionApprove, ActionRollback}
	idSet := make(map[uint]struct{})
	for _, action := range actions {
		ids, err := s.aclRepo.ListEnvironmentIDs(ctx, userID, projectID, action)
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			idSet[id] = struct{}{}
		}
	}
	if len(idSet) == 0 {
		return []uint{}, nil
	}
	result := make([]uint, 0, len(idSet))
	for id := range idSet {
		result = append(result, id)
	}
	return result, nil
}

// Get 查询
func (s *environmentManager) Get(ctx context.Context, projectID, environmentID uint) (dto.EnvironmentDTO, error) {
	environment, err := s.repo.Get(ctx, projectID, environmentID)
	if err != nil {
		return dto.EnvironmentDTO{}, err
	}
	return dto.EnvironmentDTO{
		ID:            environment.ID,
		Name:          environment.Name,
		ProjectID:     environment.ProjectID,
		ProjectName:   environment.Project.Name,
		HarborID:      environment.HarborID,
		HarborBaseURL: environment.Harbor.BaseURL,
		Description:   environment.Description,
		Creator:       environment.Creator,
		CreatedAt:     environment.CreatedAt,
		UpdatedAt:     environment.UpdatedAt,
	}, nil
}

// Create 创建
func (s *environmentManager) Create(ctx context.Context, projectID uint, req *request.CreateEnvironmentRequest) error {
	data := &model.Environment{
		Name:        req.Name,
		ProjectID:   projectID,
		HarborID:    req.HarborID,
		Description: req.Description,
		Creator:     req.Creator,
	}
	credential, err := s.credRepo.Get(ctx, req.HarborID)
	if err != nil {
		return err
	}
	// 检查凭证类型是否为 Harbor
	if credential.Type != "harbor" {
		return errors.New("credential type must be harbor")
	}
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *environmentManager) Update(ctx context.Context, projectID, environmentID uint, req *request.UpdateEnvironmentRequest) error {
	credential, err := s.credRepo.Get(ctx, *req.HarborID)
	if err != nil {
		return err
	}
	// 检查凭证类型是否为 Harbor
	if credential.Type != "harbor" {
		return errors.New("credential type must be harbor")
	}
	// 构建 Model
	data := &model.Environment{}
	data.ApplyRequest(req)
	return s.repo.Update(ctx, projectID, environmentID, data)
}

// Delete 更新
func (s *environmentManager) Delete(ctx context.Context, projectID, environmentID uint) error {
	return s.repo.Delete(ctx, projectID, environmentID)
}
