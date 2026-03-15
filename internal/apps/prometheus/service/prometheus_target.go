package service

import (
	"context"
	"valyria-backend/internal/apps/prometheus/dto"
	"valyria-backend/internal/apps/prometheus/executor"
	"valyria-backend/internal/apps/prometheus/model"
	"valyria-backend/internal/apps/prometheus/repository"
	"valyria-backend/internal/apps/prometheus/request"
	pg "valyria-backend/internal/core/pagination"
)

// TargetManager 定义服务层接口
type TargetManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.TargetDTO, pg.Pagination, error)
	Get(ctx context.Context, id int) (dto.TargetDTO, error)
	Create(ctx context.Context, groupID int, req *request.CreateTargetRequest) error
	Update(ctx context.Context, id int, req *request.UpdateTargetRequest) error
	Delete(ctx context.Context, id int) error
}

// targetManager 实现 TargetManager 接口
type targetManager struct {
	target repository.TargetRepository
	syncer executor.TargetSyncer
}

// NewTargetManager 创建新的 TargetManager 实例
func NewTargetManager(target repository.TargetRepository, syncer executor.TargetSyncer) TargetManager {
	return &targetManager{target: target, syncer: syncer}
}

// List 列表
func (s *targetManager) List(ctx context.Context, params pg.QueryParams) ([]dto.TargetDTO, pg.Pagination, error) {
	data, pagination, err := s.target.List(ctx, params)
	if err != nil {
		return nil, pg.Pagination{}, err
	}
	result := make([]dto.TargetDTO, 0, len(data))
	for _, item := range data {
		result = append(result, dto.ToTargetDTO(item))
	}
	return result, pagination, nil
}

func (s *targetManager) Get(ctx context.Context, id int) (dto.TargetDTO, error) {
	item, err := s.target.Get(ctx, id)
	if err != nil {
		return dto.TargetDTO{}, err
	}
	return dto.ToTargetDTO(item), nil
}

func (s *targetManager) Create(ctx context.Context, groupID int, req *request.CreateTargetRequest) error {
	data := &model.Target{
		GroupID:   groupID,
		IPAddress: req.IPAddress,
		Port:      req.Port,
		Labels:    req.Labels,
		Enabled:   req.Enabled,
	}
	if err := s.target.Create(ctx, data); err != nil {
		return err
	}
	return s.syncer.SyncTargetGroup(ctx, groupID)
}

func (s *targetManager) Update(ctx context.Context, id int, req *request.UpdateTargetRequest) error {
	current, err := s.target.Get(ctx, id)
	if err != nil {
		return err
	}
	data := &model.Target{}
	if req.IPAddress != nil {
		data.IPAddress = *req.IPAddress
	}
	if req.Port != nil {
		data.Port = *req.Port
	}
	if req.Labels != nil {
		data.Labels = *req.Labels
	}
	if req.Enabled != nil {
		data.Enabled = req.Enabled
	}
	if err := s.target.Update(ctx, id, data); err != nil {
		return err
	}
	return s.syncer.SyncTargetGroup(ctx, current.GroupID)
}

func (s *targetManager) Delete(ctx context.Context, id int) error {
	current, err := s.target.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.target.Delete(ctx, id); err != nil {
		return err
	}
	return s.syncer.SyncTargetGroup(ctx, current.GroupID)
}

