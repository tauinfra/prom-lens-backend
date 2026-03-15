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

// RecordManager 定义服务层接口
type RecordManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.RecordDTO, pg.Pagination, error)
	Get(ctx context.Context, id int) (dto.RecordDTO, error)
	Create(ctx context.Context, groupID int, req *request.CreateRecordRequest) error
	Update(ctx context.Context, id int, req *request.UpdateRecordRequest) error
	Delete(ctx context.Context, id int) error
}

// recordManager 实现 RecordManager 接口
type recordManager struct {
	record repository.RecordRepository
	syncer executor.RuleSyncer
}

// NewRecordManager 创建新的 RecordManager 实例
func NewRecordManager(record repository.RecordRepository, syncer executor.RuleSyncer) RecordManager {
	return &recordManager{record: record, syncer: syncer}
}

// List 列表
func (s *recordManager) List(ctx context.Context, params pg.QueryParams) ([]dto.RecordDTO, pg.Pagination, error) {
	data, pagination, err := s.record.List(ctx, params)
	if err != nil {
		return nil, pg.Pagination{}, err
	}
	result := make([]dto.RecordDTO, 0, len(data))
	for _, item := range data {
		result = append(result, dto.ToRecordDTO(item))
	}
	return result, pagination, nil
}

func (s *recordManager) Get(ctx context.Context, id int) (dto.RecordDTO, error) {
	item, err := s.record.Get(ctx, id)
	if err != nil {
		return dto.RecordDTO{}, err
	}
	return dto.ToRecordDTO(item), nil
}

func (s *recordManager) Create(ctx context.Context, groupID int, req *request.CreateRecordRequest) error {
	data := &model.Record{
		Name:    req.Name,
		Expr:    req.Expr,
		GroupID: groupID,
	}
	if err := s.record.Create(ctx, data); err != nil {
		return err
	}
	return s.syncer.SyncRuleGroup(ctx, groupID)
}

func (s *recordManager) Update(ctx context.Context, id int, req *request.UpdateRecordRequest) error {
	current, err := s.record.Get(ctx, id)
	if err != nil {
		return err
	}
	data := &model.Record{}
	if req.Name != nil {
		data.Name = *req.Name
	}
	if req.Expr != nil {
		data.Expr = *req.Expr
	}
	if err := s.record.Update(ctx, id, data); err != nil {
		return err
	}
	return s.syncer.SyncRuleGroup(ctx, current.GroupID)
}

func (s *recordManager) Delete(ctx context.Context, id int) error {
	current, err := s.record.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.record.Delete(ctx, id); err != nil {
		return err
	}
	return s.syncer.SyncRuleGroup(ctx, current.GroupID)
}

