package service

import (
	"context"
	"errors"
	"fmt"

	prom "prom-lens-backend/internal/apps/prometheus"
	"prom-lens-backend/internal/apps/prometheus/dto"
	"prom-lens-backend/internal/apps/prometheus/model"
	"prom-lens-backend/internal/apps/prometheus/repository"
	"prom-lens-backend/internal/apps/prometheus/request"
	pg "prom-lens-backend/internal/core/pagination"
)

// GroupManager 定义接口
type GroupManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.GroupDTO, pg.Pagination, error)
	Get(ctx context.Context, id int) (dto.GroupDTO, error)
	Create(ctx context.Context, req *request.CreateGroupRequest) error
	Update(ctx context.Context, id int, req *request.UpdateGroupRequest) error
	Delete(ctx context.Context, id int) error
}

// groupManager 实现 GroupManager 接口
type groupManager struct {
	group repository.GroupRepository
}

// NewGroupManager 创建新的 GroupManager 实例
func NewGroupManager(group repository.GroupRepository) GroupManager {
	return &groupManager{group: group}
}

// List 列表
func (s *groupManager) List(ctx context.Context, params pg.QueryParams) ([]dto.GroupDTO, pg.Pagination, error) {
	data, pagination, err := s.group.List(ctx, params)
	if err != nil {
		return nil, pg.Pagination{}, err
	}
	result := make([]dto.GroupDTO, 0, len(data))
	for _, item := range data {
		groupDTO, err := s.toGroupDTO(ctx, item)
		if err != nil {
			return nil, pg.Pagination{}, err
		}
		result = append(result, groupDTO)
	}
	return result, pagination, nil
}

func (s *groupManager) toGroupDTO(ctx context.Context, item model.Group) (dto.GroupDTO, error) {
	ruleCount, err := s.group.CountRules(ctx, item.ID)
	if err != nil {
		return dto.GroupDTO{}, err
	}
	recordCount, err := s.group.CountRecords(ctx, item.ID)
	if err != nil {
		return dto.GroupDTO{}, err
	}
	return dto.ToGroupDTO(item, ruleCount, recordCount), nil
}

func (s *groupManager) Get(ctx context.Context, id int) (dto.GroupDTO, error) {
	item, err := s.group.Get(ctx, id)
	if err != nil {
		return dto.GroupDTO{}, err
	}
	return s.toGroupDTO(ctx, item)
}

func (s *groupManager) Create(ctx context.Context, req *request.CreateGroupRequest) error {
	if !prom.IsAlertingRules(req.Type) && !prom.IsAlertingRecords(req.Type) {
		return fmt.Errorf("invalid group type: %s", req.Type)
	}
	data := &model.Group{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
	}
	return s.group.Create(ctx, data)
}

func (s *groupManager) Update(ctx context.Context, id int, req *request.UpdateGroupRequest) error {
	data := &model.Group{}
	if req.Name != nil {
		data.Name = *req.Name
	}
	if req.Type != nil {
		if !prom.IsAlertingRules(*req.Type) && !prom.IsAlertingRecords(*req.Type) {
			return fmt.Errorf("invalid group type: %s", *req.Type)
		}
		data.Type = *req.Type
	}
	if req.Description != nil {
		data.Description = *req.Description
	}
	return s.group.Update(ctx, id, data)
}

func (s *groupManager) Delete(ctx context.Context, id int) error {
	exists, err := s.group.HasRules(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("该分组下仍存在 Rule，禁止删除")
	}
	exists, err = s.group.HasRecords(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("该分组下仍存在 Record，禁止删除")
	}
	return s.group.Delete(ctx, id)
}
