package service

import (
	"context"
	"gorm.io/gorm"
	"valyria-backend/internal/apps/kubernetes/model"
	"valyria-backend/internal/apps/kubernetes/repository"
	pg "valyria-backend/internal/core/pagination"
)

type ClusterService interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Cluster, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Cluster, error)
	Create(ctx context.Context, data *model.Cluster) error
	Update(ctx context.Context, id int, user *model.Cluster) error
	Delete(ctx context.Context, id int) error
}

// clusterService 实现 ClusterService 接口
type clusterService struct {
	repo repository.ClusterRepository
	db   *gorm.DB
}

// NewClusterService 创建新的 ClusterService 实例
func NewClusterService(db *gorm.DB, repo repository.ClusterRepository) ClusterService {
	return &clusterService{db: db, repo: repo}
}

// List 列表
func (s *clusterService) List(ctx context.Context, params pg.QueryParams) ([]model.Cluster, pg.Pagination, error) {
	return s.repo.List(ctx, params)
}

// Get 查询
func (s *clusterService) Get(ctx context.Context, id int) (model.Cluster, error) {
	return s.repo.Get(ctx, id)
}

// Create 创建
func (s *clusterService) Create(ctx context.Context, data *model.Cluster) error {
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *clusterService) Update(ctx context.Context, id int, data *model.Cluster) error {
	return s.repo.Update(ctx, id, data)
}

// Delete 删除
func (s *clusterService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
