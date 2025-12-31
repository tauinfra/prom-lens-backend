package service

import (
	"context"
	"valyria-backend/internal/apps/auth/model"
	"valyria-backend/internal/apps/auth/repository"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type UserService interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.User, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.User, error)
	Create(ctx context.Context, data *model.User) error
	Update(ctx context.Context, id int, data *model.User) error
	Delete(ctx context.Context, id int) error
	Login(ctx context.Context, username, password string) (model.User, error)
	ChangePassword(ctx context.Context, username, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, username, newPassword string) error
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
func (s *userService) List(ctx context.Context, params pg.QueryParams) ([]model.User, pg.Pagination, error) {
	return s.repo.List(ctx, params)
}

// Get 查询
func (s *userService) Get(ctx context.Context, id int) (model.User, error) {
	return s.repo.Get(ctx, id)
}

// Create 创建
func (s *userService) Create(ctx context.Context, data *model.User) error {
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *userService) Update(ctx context.Context, id int, data *model.User) error {
	return s.repo.Update(ctx, id, data)
}

// Delete 删除
func (s *userService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// Login 登录
func (s *userService) Login(ctx context.Context, username, password string) (model.User, error) {
	return s.repo.Login(ctx, username, password)
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, username, oldPassword, newPassword string) error {
	// 更新用户密码
	return s.repo.UpdatePassword(ctx, username, oldPassword, newPassword)
}

// ResetPassword 重置密码
func (s *userService) ResetPassword(ctx context.Context, username, newPassword string) error {
	// 重置密码
	return s.repo.ResetPassword(ctx, username, newPassword)
}
