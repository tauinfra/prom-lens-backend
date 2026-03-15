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

type RoleManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.RoleDTO, pg.Pagination, error)
	Get(ctx context.Context, id uint) (dto.RoleDTO, error)
	Create(ctx context.Context, req *request.CreateRoleRequest) error
	Update(ctx context.Context, id uint, req *request.UpdateRoleRequest) error
	Delete(ctx context.Context, id uint) error
	GetRolePermissions(ctx context.Context, roleID uint) ([]dto.PermissionDTO, error)
	UpdateRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
}

// roleManager 实现 RoleManager 接口
type roleManager struct {
	repo repository.RoleRepository
	db   *gorm.DB
}

// NewRoleManager 创建新的 RoleManager 实例
func NewRoleManager(db *gorm.DB, repo repository.RoleRepository) RoleManager {
	return &roleManager{db: db, repo: repo}
}

// List 列表
func (s *roleManager) List(ctx context.Context, params pg.QueryParams) (data []dto.RoleDTO, pagination pg.Pagination, err error) {
	roles, pagination, err := s.repo.List(ctx, params)
	if err != nil {
		return data, pagination, err
	}
	for _, role := range roles {
		data = append(data, dto.RoleDTO{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
			Creator:     role.Creator,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}
	return data, pagination, nil
}

// Get 查询
func (s *roleManager) Get(ctx context.Context, id uint) (dto.RoleDTO, error) {
	role, err := s.repo.Get(ctx, id)
	if err != nil {
		return dto.RoleDTO{}, err
	}
	return dto.RoleDTO{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		Creator:     role.Creator,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}

// Create 创建
func (s *roleManager) Create(ctx context.Context, req *request.CreateRoleRequest) error {
	// 构建 Model
	data := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Creator:     req.Creator,
	}
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *roleManager) Update(ctx context.Context, id uint, req *request.UpdateRoleRequest) error {
	// 构建 Model
	role := &model.Role{}
	role.ApplyRequest(req)
	// 更新 Role
	return s.repo.Update(ctx, id, role)
}

// Delete 删除
func (s *roleManager) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// GetRolePermissions  查看角色权限
func (s *roleManager) GetRolePermissions(ctx context.Context, roleID uint) (data []dto.PermissionDTO, err error) {
	permissions, err := s.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}
	for _, perm := range permissions {
		data = append(data, dto.PermissionDTO{
			ID:        perm.ID,
			Name:      perm.Name,
			Code:      perm.Code,
			Creator:   perm.Creator,
			CreatedAt: perm.CreatedAt,
			UpdatedAt: perm.UpdatedAt,
		})
	}
	return data, nil
}

// UpdateRolePermissions 更新角色权限
func (s *roleManager) UpdateRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return s.repo.UpdateRolePermissions(ctx, roleID, permissionIDs)
}
