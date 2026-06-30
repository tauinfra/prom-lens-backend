package service

import (
	"context"
	"fmt"

	prom "prom-lens-backend/internal/apps/prometheus"
	"prom-lens-backend/internal/apps/prometheus/dto"
	"prom-lens-backend/internal/apps/prometheus/executor"
	"prom-lens-backend/internal/apps/prometheus/model"
	"prom-lens-backend/internal/apps/prometheus/repository"
	"prom-lens-backend/internal/apps/prometheus/request"
	pg "prom-lens-backend/internal/core/pagination"
)

// RuleManager 定义服务层接口
type RuleManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.RuleDTO, pg.Pagination, error)
	Get(ctx context.Context, id int) (dto.RuleDTO, error)
	Create(ctx context.Context, groupID int, req *request.CreateRuleRequest) error
	Update(ctx context.Context, id int, req *request.UpdateRuleRequest) error
	Delete(ctx context.Context, id int) error
}

// ruleManager 实现 RuleManager 接口
type ruleManager struct {
	group  repository.GroupRepository
	rule   repository.RuleRepository
	syncer executor.RuleSyncer
}

// NewRuleManager 创建新的 RuleManager 实例
func NewRuleManager(group repository.GroupRepository, rule repository.RuleRepository, syncer executor.RuleSyncer) RuleManager {
	return &ruleManager{group: group, rule: rule, syncer: syncer}
}

// List 列表
func (s *ruleManager) List(ctx context.Context, params pg.QueryParams) ([]dto.RuleDTO, pg.Pagination, error) {
	data, pagination, err := s.rule.List(ctx, params)
	if err != nil {
		return nil, pg.Pagination{}, err
	}
	result := make([]dto.RuleDTO, 0, len(data))
	for _, item := range data {
		result = append(result, dto.ToRuleDTO(item))
	}
	return result, pagination, nil
}

func (s *ruleManager) Get(ctx context.Context, id int) (dto.RuleDTO, error) {
	item, err := s.rule.Get(ctx, id)
	if err != nil {
		return dto.RuleDTO{}, err
	}
	return dto.ToRuleDTO(item), nil
}

func (s *ruleManager) Create(ctx context.Context, groupID int, req *request.CreateRuleRequest) error {
	group, err := s.group.Get(ctx, groupID)
	if err != nil {
		return err
	}
	if !prom.IsAlertingRules(group.Type) {
		return fmt.Errorf("group type is %q, only %q groups can contain rules", group.Type, prom.GroupTypeAlertingRules)
	}
	if err := prom.ValidateExtraAnnotations(req.ExtraAnnotations); err != nil {
		return err
	}
	data := &model.Rule{
		Name:              req.Name,
		GroupID:           groupID,
		Summary:           req.Summary,
		Description:       req.Description,
		Expr:              req.Expr,
		For:               req.For,
		Labels:            req.Labels,
		ExtraAnnotations:  prom.NormalizeExtraAnnotations(req.ExtraAnnotations),
		Status:            req.Status,
	}
	if err := s.rule.Create(ctx, data); err != nil {
		return err
	}
	return s.syncer.SyncRuleGroup(ctx, groupID)
}

func (s *ruleManager) Update(ctx context.Context, id int, req *request.UpdateRuleRequest) error {
	current, err := s.rule.Get(ctx, id)
	if err != nil {
		return err
	}
	data := &model.Rule{}
	if req.Name != nil {
		data.Name = *req.Name
	}
	if req.Summary != nil {
		data.Summary = *req.Summary
	}
	if req.Description != nil {
		data.Description = *req.Description
	}
	if req.Expr != nil {
		data.Expr = *req.Expr
	}
	if req.For != nil {
		data.For = *req.For
	}
	if req.Labels != nil {
		data.Labels = *req.Labels
	}
	if req.ExtraAnnotations != nil {
		if err := prom.ValidateExtraAnnotations(*req.ExtraAnnotations); err != nil {
			return err
		}
		data.ExtraAnnotations = prom.NormalizeExtraAnnotations(*req.ExtraAnnotations)
	}
	if req.Status != nil {
		data.Status = req.Status
	}
	if err := s.rule.Update(ctx, id, data); err != nil {
		return err
	}
	return s.syncer.SyncRuleGroup(ctx, current.GroupID)
}

func (s *ruleManager) Delete(ctx context.Context, id int) error {
	current, err := s.rule.Get(ctx, id)
	if err != nil {
		return err
	}
	if err = s.rule.Delete(ctx, id); err != nil {
		return err
	}
	return s.syncer.SyncRuleGroup(ctx, current.GroupID)
}
