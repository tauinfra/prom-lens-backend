package service

import (
	"context"
	"errors"
	"time"
	"prom-lens-backend/internal/apps/authn/model"
	"prom-lens-backend/internal/apps/authn/repository"
	"prom-lens-backend/internal/apps/authn/request"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthManager interface {
	Login(ctx context.Context, username, password string) (model.User, error)
	GetMe(ctx context.Context, id uint) (model.User, error)
	UpdateLastLogin(ctx context.Context, userID uint) error
	ChangePassword(ctx context.Context, id uint, req *request.ChangePasswordRequest) error
}

type authManager struct {
	repo repository.AuthRepository
	db   *gorm.DB
}

func NewAuthManager(db *gorm.DB, repo repository.AuthRepository) AuthManager {
	return &authManager{db: db, repo: repo}
}

func (s *authManager) Login(ctx context.Context, username, password string) (model.User, error) {
	return s.repo.Login(ctx, username, password)
}

func (s *authManager) GetMe(ctx context.Context, id uint) (model.User, error) {
	return s.repo.GetMe(ctx, id)
}

func (s *authManager) UpdateLastLogin(ctx context.Context, userID uint) error {
	return s.repo.UpdateLastLogin(ctx, userID, time.Now())
}

func (s *authManager) ChangePassword(ctx context.Context, id uint, req *request.ChangePasswordRequest) error {
	user, err := s.repo.GetMe(ctx, id)
	if err != nil {
		return err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("旧密码错误")
	}
	newPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.ChangeMyPassword(ctx, id, string(newPassword))
}
