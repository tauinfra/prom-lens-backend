package service

import (
	"context"
	"errors"
	"prom-lens-backend/internal/apps/prometheus/dto"
	"prom-lens-backend/internal/apps/prometheus/model"
	"prom-lens-backend/internal/apps/prometheus/repository"
	"prom-lens-backend/internal/apps/prometheus/request"
	pg "prom-lens-backend/internal/core/pagination"
)

// TargetGroupManager 定义接口
type TargetGroupManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.TargetGroupDTO, pg.Pagination, error)
	Get(ctx context.Context, id int) (dto.TargetGroupDTO, error)
	Create(ctx context.Context, req *request.CreateTargetGroupRequest) error
	Update(ctx context.Context, id int, req *request.UpdateTargetGroupRequest) error
	Delete(ctx context.Context, id int) error
}

// targetGroupManager 实现 TargetGroupManager 接口
type targetGroupManager struct {
	group repository.TargetGroupRepository
}

// NewTargetGroupManager 创建新的 TargetGroupManager 实例
func NewTargetGroupManager(group repository.TargetGroupRepository) TargetGroupManager {
	return &targetGroupManager{group: group}
}

// List 列表
func (s *targetGroupManager) List(ctx context.Context, params pg.QueryParams) ([]dto.TargetGroupDTO, pg.Pagination, error) {
	data, pagination, err := s.group.List(ctx, params)
	if err != nil {
		return nil, pg.Pagination{}, err
	}
	result := make([]dto.TargetGroupDTO, 0, len(data))
	for _, item := range data {
		result = append(result, dto.ToTargetGroupDTO(item))
	}
	return result, pagination, nil
}

func (s *targetGroupManager) Get(ctx context.Context, id int) (dto.TargetGroupDTO, error) {
	item, err := s.group.Get(ctx, id)
	if err != nil {
		return dto.TargetGroupDTO{}, err
	}
	return dto.ToTargetGroupDTO(item), nil
}

func (s *targetGroupManager) Create(ctx context.Context, req *request.CreateTargetGroupRequest) error {
	data := &model.TargetGroup{
		Name:        req.Name,
		Description: req.Description,
		Labels:      req.Labels,
	}
	return s.group.Create(ctx, data)
}

func (s *targetGroupManager) Update(ctx context.Context, id int, req *request.UpdateTargetGroupRequest) error {
	data := &model.TargetGroup{}
	if req.Name != nil {
		data.Name = *req.Name
	}
	if req.Description != nil {
		data.Description = *req.Description
	}
	if req.Labels != nil {
		data.Labels = *req.Labels
	}
	return s.group.Update(ctx, id, data)
}

func (s *targetGroupManager) Delete(ctx context.Context, id int) error {
	exists, err := s.group.HasTargets(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("该分组下仍存在 Target 实例，请先清空后再删除")
	}
	return s.group.Delete(ctx, id)
}
