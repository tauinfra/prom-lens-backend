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

type ProjectManager interface {
	List(ctx context.Context, userID uint, isAdmin bool, params pg.QueryParams) ([]dto.ProjectDTO, pg.Pagination, error)
	Get(ctx context.Context, id uint) (dto.ProjectDTO, error)
	Create(ctx context.Context, req *request.CreateProjectRequest, creator string) error
	Update(ctx context.Context, id uint, req *request.UpdateProjectRequest) error
	Delete(ctx context.Context, id uint) error
}

// projectManager 实现 ProjectManager 接口
type projectManager struct {
	repo      repository.ProjectRepository
	aclRepo   repository.PipelineACLRepository
	db        *gorm.DB
}

// NewProjectManager 创建新的 ProjectManager 实例
func NewProjectManager(repo repository.ProjectRepository, aclRepo repository.PipelineACLRepository, db *gorm.DB) ProjectManager {
	return &projectManager{repo: repo, aclRepo: aclRepo, db: db}
}

// List 列表
func (s *projectManager) List(ctx context.Context, userID uint, isAdmin bool, params pg.QueryParams) (data []dto.ProjectDTO, pagination pg.Pagination, err error) {
	var ids []uint
	if !isAdmin {
		if userID == 0 {
			return []dto.ProjectDTO{}, pagination, errors.New("no permission")
		}
		if ids, err = s.listViewableProjectIDs(ctx, userID); err != nil {
			return
		}
		if len(ids) == 0 {
			return []dto.ProjectDTO{}, pg.Pagination{Page: params.Page, Size: params.Size, Total: 0}, nil
		}
		params.Conditions = []pg.Condition{
			{Field: "id", Operator: "IN", Value: ids},
		}
	}
	list, pagination, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, pg.Pagination{}, err
	}
	data = make([]dto.ProjectDTO, 0, len(list))
	for _, p := range list {
		data = append(data, projectToDTO(p))
	}
	return data, pagination, nil
}

func (s *projectManager) listViewableProjectIDs(ctx context.Context, userID uint) ([]uint, error) {
	actions := []string{ActionView, ActionDeploy, ActionApprove, ActionRollback}
	idSet := make(map[uint]struct{})
	for _, action := range actions {
		ids, err := s.aclRepo.ListProjectIDs(ctx, userID, action)
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
func (s *projectManager) Get(ctx context.Context, id uint) (dto.ProjectDTO, error) {
	p, err := s.repo.Get(ctx, id)
	if err != nil {
		return dto.ProjectDTO{}, err
	}
	return projectToDTO(p), nil
}

func projectToDTO(p model.Project) dto.ProjectDTO {
	return dto.ProjectDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Creator:     p.Creator,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// Create 创建
func (s *projectManager) Create(ctx context.Context, req *request.CreateProjectRequest, creator string) error {
	data := &model.Project{
		Name:        req.Name,
		Description: req.Description,
		Creator:     creator,
	}
	return s.repo.Create(ctx, data)
}

// Update 更新（DTO 模式：仅更新请求中传入的字段）
func (s *projectManager) Update(ctx context.Context, id uint, req *request.UpdateProjectRequest) error {
	data := &model.Project{}
	var fields []string
	if req.Name != nil {
		data.Name = *req.Name
		fields = append(fields, "name")
	}
	if req.Description != nil {
		data.Description = *req.Description
		fields = append(fields, "description")
	}
	return s.repo.Update(ctx, id, data, fields)
}

// Delete 更新
func (s *projectManager) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
