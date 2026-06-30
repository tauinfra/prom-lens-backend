package service

import (
	"context"
	"errors"
	"fmt"

	"prom-lens-backend/internal/apps/authn/dto"
	"prom-lens-backend/internal/apps/authn/model"
	"prom-lens-backend/internal/apps/authn/repository"
	"prom-lens-backend/internal/apps/authn/request"
	pg "prom-lens-backend/internal/core/pagination"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.UserDTO, pg.Pagination, error)
	Get(ctx context.Context, id uint) (dto.UserDTO, error)
	Create(ctx context.Context, req *request.CreateUserRequest) error
	Update(ctx context.Context, id uint, req *request.UpdateUserRequest) error
	Delete(ctx context.Context, operatorID uint, id uint) error
	ResetPassword(ctx context.Context, id uint, req *request.ResetPasswordRequest) error
}

type userManager struct {
	repo repository.UserRepository
	db   *gorm.DB
}

func NewUserManager(db *gorm.DB, repo repository.UserRepository) UserManager {
	return &userManager{db: db, repo: repo}
}

func (s *userManager) List(ctx context.Context, params pg.QueryParams) ([]dto.UserDTO, pg.Pagination, error) {
	users, pagination, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, pagination, err
	}
	data := make([]dto.UserDTO, 0, len(users))
	for _, user := range users {
		data = append(data, toUserDTO(user))
	}
	return data, pagination, nil
}

func (s *userManager) Get(ctx context.Context, id uint) (dto.UserDTO, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return dto.UserDTO{}, err
	}
	return toUserDTO(user), nil
}

func (s *userManager) Create(ctx context.Context, req *request.CreateUserRequest) error {
	var dn *string
	if req.DN != "" {
		dn = &req.DN
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	isSuperuser := false
	if req.IsSuperuser != nil {
		isSuperuser = *req.IsSuperuser
	}
	isLdap := false
	if req.IsLdap != nil {
		isLdap = *req.IsLdap
	}
	data := &model.User{
		Username:    req.Username,
		Password:    req.Password,
		Nickname:    req.Nickname,
		Email:       req.Email,
		Phone:       req.Phone,
		IsActive:    &isActive,
		IsSuperuser: &isSuperuser,
		IsLdap:      &isLdap,
		DN:          dn,
		Creator:     req.Creator,
	}
	return s.repo.Create(ctx, data)
}

func (s *userManager) Update(ctx context.Context, id uint, req *request.UpdateUserRequest) error {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if req.IsSuperuser != nil && !*req.IsSuperuser && user.IsSuperuser != nil && *user.IsSuperuser {
		count, err := s.repo.CountSuperusers(ctx)
		if err != nil {
			return err
		}
		if count <= 1 {
			return errors.New("不能取消唯一超级管理员权限")
		}
	}
	update := &model.User{}
	update.ApplyRequest(req)
	return s.repo.Update(ctx, id, update)
}

func (s *userManager) Delete(ctx context.Context, operatorID uint, id uint) error {
	if operatorID == id {
		return errors.New("不能删除当前登录用户")
	}
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if user.IsSuperuser != nil && *user.IsSuperuser {
		count, err := s.repo.CountSuperusers(ctx)
		if err != nil {
			return err
		}
		if count <= 1 {
			return errors.New("不能删除唯一超级管理员")
		}
	}
	return s.repo.Delete(ctx, id)
}

func (s *userManager) ResetPassword(ctx context.Context, id uint, req *request.ResetPasswordRequest) error {
	if _, err := s.repo.Get(ctx, id); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}
	return s.repo.ResetPassword(ctx, id, string(hash))
}

func toUserDTO(user model.User) dto.UserDTO {
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
		LastLogin:   user.LastLogin,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
