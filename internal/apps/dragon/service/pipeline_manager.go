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

type PipelineManager interface {
	List(ctx context.Context, projectID, environmentID uint, params pg.QueryParams) ([]dto.PipelineDTO, pg.Pagination, error)
	Get(ctx context.Context, projectID, environmentID, pipelineID uint) (dto.PipelineDTO, error)
	Create(ctx context.Context, projectID, environmentID uint, req *request.CreatePipelineRequest) error
	Update(ctx context.Context, projectID, environmentID, pipelineID uint, req *request.UpdatePipelineRequest) error
	Delete(ctx context.Context, projectID, environmentID, pipelineID uint) error
}

// pipelineManager 实现 PipelineManager 接口
type pipelineManager struct {
	repo    repository.PipelineRepository
	aclRepo repository.PipelineACLRepository
	db      *gorm.DB
}

// NewPipelineManager 创建新的 PipelineManager 实例
func NewPipelineManager(
	repo repository.PipelineRepository,
	aclRepo repository.PipelineACLRepository,
	db *gorm.DB,
) PipelineManager {
	return &pipelineManager{repo: repo, aclRepo: aclRepo, db: db}
}

// List 列表
func (s *pipelineManager) List(ctx context.Context, projectID, environmentID uint, params pg.QueryParams) (data []dto.PipelineDTO, pagination pg.Pagination, err error) {
	userID := GetUID(ctx)
	if !IsSuperuser(ctx) {
		if userID == 0 {
			return data, pagination, errors.New("no permission")
		}
		ids, err := s.listViewablePipelineIDs(ctx, userID)
		if err != nil {
			return data, pagination, err
		}
		if len(ids) == 0 {
			return []dto.PipelineDTO{}, pg.Pagination{Page: params.Page, Size: params.Size, Total: 0}, nil
		}
		params.Conditions = append(params.Conditions, pg.Condition{
			Field:    "valyria_dragon_pipeline.id",
			Operator: "IN",
			Value:    ids,
		})
	}
	pipelines, pagination, err := s.repo.List(ctx, projectID, environmentID, params)
	if err != nil {
		return data, pagination, err
	}
	for _, pipeline := range pipelines {
		data = append(data, dto.PipelineDTO{
			ID:              pipeline.ID,
			Name:            pipeline.Name,
			EnvironmentID:   pipeline.EnvironmentID,
			EnvironmentName: pipeline.Environment.Name,
			ProjectID:       projectID,
			ProjectName:     pipeline.Environment.Project.Name,
			HarborBaseURL:   pipeline.Environment.Harbor.BaseURL,
			GitlabID:        pipeline.GitlabID,
			GitlabGroupID:   pipeline.GitlabGroupID,
			GitlabProjectID: pipeline.GitlabProjectID,
			GitlabBaseURL:   pipeline.Gitlab.BaseURL,
			IsApproval:      pipeline.IsApproval,
			Task:            pipeline.Task,
			Creator:         pipeline.Creator,
			CreatedAt:       pipeline.CreatedAt,
			UpdatedAt:       pipeline.UpdatedAt,
		})
	}
	return data, pagination, nil
}

func (s *pipelineManager) listViewablePipelineIDs(ctx context.Context, userID uint) ([]uint, error) {
	actions := []string{ActionView, ActionDeploy, ActionApprove, ActionRollback}
	idSet := make(map[uint]struct{})
	for _, action := range actions {
		ids, err := s.aclRepo.ListPipelineIDs(ctx, userID, action)
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
func (s *pipelineManager) Get(ctx context.Context, projectID, environmentID, pipelineID uint) (dto.PipelineDTO, error) {
	if !IsSuperuser(ctx) {
		userID := GetUID(ctx)
		if userID == 0 {
			return dto.PipelineDTO{}, errors.New("no permission")
		}
		ok, err := s.hasViewPermission(ctx, userID, pipelineID)
		if err != nil {
			return dto.PipelineDTO{}, err
		}
		if !ok {
			return dto.PipelineDTO{}, errors.New("no permission")
		}
	}
	pipeline, err := s.repo.Get(ctx, projectID, environmentID, pipelineID)
	if err != nil {
		return dto.PipelineDTO{}, err
	}
	environmentName := ""
	projectIDVal := uint(0)
	projectName := ""
	harborBaseURL := ""
	if pipeline.Environment != nil {
		environmentName = pipeline.Environment.Name
		projectIDVal = pipeline.Environment.ProjectID
		if pipeline.Environment.Project != nil {
			projectName = pipeline.Environment.Project.Name
		}
		if pipeline.Environment.Harbor != nil {
			harborBaseURL = pipeline.Environment.Harbor.BaseURL
		}
	}
	gitlabBaseURL := ""
	if pipeline.Gitlab != nil {
		gitlabBaseURL = pipeline.Gitlab.BaseURL
	}
	return dto.PipelineDTO{
		ID:              pipeline.ID,
		Name:            pipeline.Name,
		EnvironmentID:   pipeline.EnvironmentID,
		EnvironmentName: environmentName,
		ProjectID:       projectIDVal,
		ProjectName:     projectName,
		HarborBaseURL:   harborBaseURL,
		GitlabID:        pipeline.GitlabID,
		GitlabGroupID:   pipeline.GitlabGroupID,
		GitlabProjectID: pipeline.GitlabProjectID,
		GitlabBaseURL:   gitlabBaseURL,
		IsApproval:      pipeline.IsApproval,
		Task:            pipeline.Task,
		Creator:         pipeline.Creator,
		CreatedAt:       pipeline.CreatedAt,
		UpdatedAt:       pipeline.UpdatedAt,
	}, nil
}

func (s *pipelineManager) hasViewPermission(ctx context.Context, userID, pipelineID uint) (bool, error) {
	actions := []string{ActionView, ActionDeploy, ActionApprove, ActionRollback}
	for _, action := range actions {
		ok, err := s.aclRepo.HasUserAction(ctx, userID, pipelineID, action)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// Create 创建
func (s *pipelineManager) Create(ctx context.Context, projectID, environmentID uint, req *request.CreatePipelineRequest) error {
	// 构建 Model
	data := &model.Pipeline{
		Name:            req.Name,
		EnvironmentID:   environmentID,
		GitlabID:        req.GitlabID,
		GitlabGroupID:   req.GitlabGroupID,
		GitlabProjectID: req.GitlabProjectID,
		IsApproval:      req.IsApproval,
		Task:            req.Task,
		Script:          req.Script,
		Creator:         req.Creator,
	}
	// 组装 DTO
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *pipelineManager) Update(ctx context.Context, projectID, environmentID, pipelineID uint, req *request.UpdatePipelineRequest) error {
	// 构建 Model
	pipeline := &model.Pipeline{}
	pipeline.ApplyUpdate(req)
	return s.repo.Update(ctx, projectID, environmentID, pipelineID, pipeline)
}

// Delete 更新
func (s *pipelineManager) Delete(ctx context.Context, projectID, environmentID, pipelineID uint) error {
	return s.repo.Delete(ctx, projectID, environmentID, pipelineID)
}
