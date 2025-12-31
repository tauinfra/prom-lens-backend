package service

import (
	"context"
	"valyria-backend/internal/apps/kingsguard/model"
	"valyria-backend/internal/apps/kingsguard/repository"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// PolicyService 定义接口
type PolicyService interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.CasbinPolicy, pg.Pagination, error)
	Create(data model.CasbinPolicy) (bool, error)
	Delete(ctx context.Context, id int) (bool, error)
}

// policyRepository 实现了 PolicyRepository 接口
type policyService struct {
	repo repository.PolicyRepository
	db   *gorm.DB
}

// NewPolicyService 创建新的 PolicyService 实例
func NewPolicyService(repo repository.PolicyRepository, db *gorm.DB) PolicyService {
	return &policyService{repo: repo, db: db}
}

// List 列表
func (s *policyService) List(ctx context.Context, params pg.QueryParams) ([]model.CasbinPolicy, pg.Pagination, error) {
	return s.repo.List(ctx, params)
}

func (s *policyService) Create(data model.CasbinPolicy) (bool, error) {
	return s.repo.Create(data)
}

func (s *policyService) Delete(ctx context.Context, id int) (bool, error) {
	return s.repo.Delete(ctx, id)
}
