package service

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/model"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/request"
	pg "valyria-backend/internal/core/pagination"
)

type ClusterManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.ClusterDTO, pg.Pagination, error)
	Get(ctx context.Context, id uint) (dto.ClusterDTO, error)
	Create(ctx context.Context, req *request.CreateClusterRequest, creator string) error
	Update(ctx context.Context, id uint, req *request.UpdateClusterRequest) error
	UpdateToken(ctx context.Context, id uint, req *request.UpdateClusterTokenRequest) error
	Delete(ctx context.Context, id uint) error
}

// clusterManager 实现 ClusterManager 接口
type clusterManager struct {
	repo repository.ClusterRepository
	db   *gorm.DB
}

// NewClusterManager 创建新的 ClusterManager 实例
func NewClusterManager(db *gorm.DB, repo repository.ClusterRepository) ClusterManager {
	return &clusterManager{db: db, repo: repo}
}

// List 列表
func (s *clusterManager) List(ctx context.Context, params pg.QueryParams) ([]dto.ClusterDTO, pg.Pagination, error) {
	clusters, pagination, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, pagination, err
	}

	result := make([]dto.ClusterDTO, 0, len(clusters))
	for _, cluster := range clusters {
		result = append(result, dto.ClusterDTO{
			ID:          cluster.ID,
			Name:        cluster.Name,
			Host:        cluster.Host,
			Version:     cluster.Version,
			Description: cluster.Description,
			CreatedAt:   cluster.CreatedAt,
			UpdatedAt:   cluster.UpdatedAt,
		})
	}
	return result, pagination, nil
}

// Get 查询
func (s *clusterManager) Get(ctx context.Context, id uint) (dto.ClusterDTO, error) {
	cluster, err := s.repo.Get(ctx, id)
	if err != nil {
		return dto.ClusterDTO{}, err
	}
	return dto.ClusterDTO{
		ID:          cluster.ID,
		Name:        cluster.Name,
		Host:        cluster.Host,
		Version:     cluster.Version,
		Description: cluster.Description,
		CreatedAt:   cluster.CreatedAt,
		UpdatedAt:   cluster.UpdatedAt,
	}, nil
}

// Create 创建
func (s *clusterManager) Create(ctx context.Context, req *request.CreateClusterRequest, creator string) error {
	// 验证 token 必填
	if req.Token == "" {
		return errors.New("token 字段为必填项")
	}

	// 转换为 model
	data := &model.Cluster{
		Name:        req.Name,
		Host:        req.Host,
		Token:       req.Token,
		Version:     req.Version,
		Description: req.Description,
		Creator:     creator,
	}
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *clusterManager) Update(ctx context.Context, id uint, req *request.UpdateClusterRequest) error {
	// 转换为 model，只更新提供的字段
	data := &model.Cluster{}
	if req.Name != nil {
		data.Name = *req.Name
	}
	if req.Host != nil {
		data.Host = *req.Host
	}
	if req.Version != nil {
		data.Version = *req.Version
	}
	if req.Description != nil {
		data.Description = *req.Description
	}
	// Token 字段不包含在更新中，禁止修改
	return s.repo.Update(ctx, id, data)
}

// UpdateToken 更新集群 token（明文输入，仓储层负责加密）
func (s *clusterManager) UpdateToken(ctx context.Context, id uint, req *request.UpdateClusterTokenRequest) error {
	if req.Token == "" {
		return errors.New("token 字段为必填项")
	}
	return s.repo.UpdateToken(ctx, id, req.Token)
}

// Delete 删除
func (s *clusterManager) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
