package service

import (
	"context"
	"valyria-backend/internal/apps/authn/dto"
	"valyria-backend/internal/apps/authn/model"
	"valyria-backend/internal/apps/authn/repository"
	"valyria-backend/internal/apps/authn/request"
	pg "valyria-backend/internal/core/pagination"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.UserDTO, pg.Pagination, error)
	Get(ctx context.Context, id uint) (dto.UserDTO, error)
	Create(ctx context.Context, req *request.CreateUserRequest) error
	Update(ctx context.Context, id uint, req *request.UpdateUserRequest) error
	Delete(ctx context.Context, id uint) error
	ResetPassword(ctx context.Context, id uint, req *request.ResetPasswordRequest) error
	GetUserRoles(ctx context.Context, id uint) (data []dto.RoleDTO, err error)
	UpdateUserRoles(ctx context.Context, id uint, roleIDs []uint) error
	GetUserMenus(ctx context.Context, id uint) (data []dto.MenuDTO, err error)
	UpdateUserMenus(ctx context.Context, id uint, menuIDs []uint) error
}

// userService 实现 UserService 接口
type userService struct {
	repo repository.UserRepository
	db   *gorm.DB
}

// NewUserService 创建新的 UserService 实例
func NewUserService(db *gorm.DB, repo repository.UserRepository) UserService {
	return &userService{db: db, repo: repo}
}

// List 列表
func (s *userService) List(ctx context.Context, params pg.QueryParams) (data []dto.UserDTO, pagination pg.Pagination, err error) {
	users, pagination, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, pagination, err
	}
	for _, user := range users {
		dn := ""
		if user.DN != nil {
			dn = *user.DN
		}
		data = append(data, dto.UserDTO{
			ID:          user.ID,
			Username:    user.Username,
			Nickname:    user.Nickname,
			Email:       user.Email,
			Phone:       user.Phone,
			IsActive:    user.IsActive,
			IsSuperuser: user.IsSuperuser,
			IsLdap:      user.IsLdap,
			DN:          dn,
			Creator:     user.Creator,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
	}
	return data, pagination, nil
}

// Get 查询
func (s *userService) Get(ctx context.Context, id uint) (dto.UserDTO, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return dto.UserDTO{}, err
	}
	dn := ""
	if user.DN != nil {
		dn = *user.DN
	}
	return dto.UserDTO{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Email:       user.Email,
		Phone:       user.Phone,
		IsActive:    user.IsActive,
		IsSuperuser: user.IsSuperuser,
		IsLdap:      user.IsLdap,
		DN:          dn,
		Creator:     user.Creator,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

// Create 创建
func (s *userService) Create(ctx context.Context, req *request.CreateUserRequest) error {
	var dn *string
	if req.DN != "" {
		dn = &req.DN
	}
	// 构建 Model
	data := &model.User{
		Username:    req.Username,
		Nickname:    req.Nickname,
		Email:       req.Email,
		Phone:       req.Phone,
		IsActive:    req.IsActive,
		IsSuperuser: req.IsSuperuser,
		IsLdap:      req.IsLdap,
		DN:          dn,
		Creator:     req.Creator,
	}
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *userService) Update(ctx context.Context, id uint, req *request.UpdateUserRequest) error {
	// 构建 Model
	user := &model.User{}
	user.ApplyRequest(req)
	// 更新用户
	return s.repo.Update(ctx, id, user)
}

// Delete 删除
func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// ResetPassword 重置密码
func (s *userService) ResetPassword(ctx context.Context, id uint, req *request.ResetPasswordRequest) error {
	// 重置密码
	// 修改新秘密
	newPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	return s.repo.ResetPassword(ctx, id, string(newPassword))
}

// GetUserRoles  查看用户角色
func (s *userService) GetUserRoles(ctx context.Context, id uint) (data []dto.RoleDTO, err error) {
	roles, err := s.repo.GetUserRoles(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		data = append(data, dto.RoleDTO{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Creator:     role.Creator,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}
	return data, nil
}

// UpdateUserRoles 更新用户角色
func (s *userService) UpdateUserRoles(ctx context.Context, id uint, roleIDs []uint) error {
	return s.repo.UpdateUserRoles(ctx, id, roleIDs)
}

// GetUserMenus 查看用户菜单
func (s *userService) GetUserMenus(ctx context.Context, id uint) (data []dto.MenuDTO, err error) {
	menus, err := s.repo.GetUserMenus(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, m := range menus {
		data = append(data, dto.MenuDTO{
			ID:         m.ID,
			ParentID:   m.ParentID,
			Name:       m.Name,
			Path:       m.Path,
			Component:  m.Component,
			Title:      m.Title,
			Icon:       m.Icon,
			Rank:       m.Rank,
			ShowParent: m.ShowParent,
			ShowLink:   m.ShowLink,
			Creator:    m.Creator,
			CreatedAt:  m.CreatedAt,
			UpdatedAt:  m.UpdatedAt,
		})
	}
	return data, nil
}

// UpdateUserMenus 更新用户菜单
func (s *userService) UpdateUserMenus(ctx context.Context, id uint, menuIDs []uint) error {
	return s.repo.UpdateUserMenus(ctx, id, menuIDs)
}
