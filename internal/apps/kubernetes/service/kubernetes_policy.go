package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/model"
	"valyria-backend/internal/apps/kubernetes/repository"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// PolicyManager 定义接口
type PolicyManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Policy, pg.Pagination, error)
	Create(data model.Policy) (bool, error)
	Delete(ctx context.Context, id int) (bool, error)
}

// policyRepository 实现了 PolicyRepository 接口
type policyManager struct {
	repo repository.PolicyRepository
	db   *gorm.DB
}

// NewPolicyManager 创建新的 PolicyManager 实例
func NewPolicyManager(repo repository.PolicyRepository, db *gorm.DB) PolicyManager {
	return &policyManager{repo: repo, db: db}
}

// List 列表
func (s *policyManager) List(ctx context.Context, params pg.QueryParams) ([]model.Policy, pg.Pagination, error) {
	return s.repo.List(ctx, params)
}

func (s *policyManager) Create(data model.Policy) (bool, error) {
	return s.repo.Create(data)
}

func (s *policyManager) Delete(ctx context.Context, id int) (bool, error) {
	return s.repo.Delete(ctx, id)
}
