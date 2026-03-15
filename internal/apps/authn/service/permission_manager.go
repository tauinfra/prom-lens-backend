package service

import (
	"context"
	"valyria-backend/internal/apps/authn/dto"
	"valyria-backend/internal/apps/authn/model"
	"valyria-backend/internal/apps/authn/repository"
	"valyria-backend/internal/apps/authn/request"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type PermissionManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.PermissionDTO, pg.Pagination, error)
	Get(ctx context.Context, id uint) (dto.PermissionDTO, error)
	Create(ctx context.Context, req *request.CreatePermissionRequest) error
	Update(ctx context.Context, id uint, req *request.UpdatePermissionRequest) error
	Delete(ctx context.Context, id uint) error
}

// permissionManager 实现 PermissionManager 接口
type permissionManager struct {
	repo repository.PermissionRepository
	db   *gorm.DB
}

// NewPermissionManager 创建新的 PermissionManager 实例
func NewPermissionManager(db *gorm.DB, repo repository.PermissionRepository) PermissionManager {
	return &permissionManager{db: db, repo: repo}
}

// List 列表
func (s *permissionManager) List(ctx context.Context, params pg.QueryParams) (data []dto.PermissionDTO, pagination pg.Pagination, err error) {
	permissions, pagination, err := s.repo.List(ctx, params)
	if err != nil {
		return data, pagination, err
	}
	for _, permission := range permissions {
		data = append(data, dto.PermissionDTO{
			ID:        permission.ID,
			Name:      permission.Name,
			Code:      permission.Code,
			Creator:   permission.Creator,
			CreatedAt: permission.CreatedAt,
			UpdatedAt: permission.UpdatedAt,
		})
	}
	return data, pagination, nil
}

// Get 查询
func (s *permissionManager) Get(ctx context.Context, id uint) (dto.PermissionDTO, error) {
	permission, err := s.repo.Get(ctx, id)
	if err != nil {
		return dto.PermissionDTO{}, err
	}
	return dto.PermissionDTO{
		ID:        permission.ID,
		Name:      permission.Name,
		Code:      permission.Code,
		Creator:   permission.Creator,
		CreatedAt: permission.CreatedAt,
		UpdatedAt: permission.UpdatedAt,
	}, nil
}

// Create 创建
func (s *permissionManager) Create(ctx context.Context, req *request.CreatePermissionRequest) error {
	// 构建 Model
	data := &model.Permission{
		Name:    req.Name,
		Code:    req.Code,
		Creator: req.Creator,
	}
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *permissionManager) Update(ctx context.Context, id uint, req *request.UpdatePermissionRequest) error {
	// 构建 Model
	permission := &model.Permission{}
	permission.ApplyRequest(req)
	// 更新 Permission
	return s.repo.Update(ctx, id, permission)
}

// Delete 删除
func (s *permissionManager) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
