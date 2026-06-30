package repository

import (
	"context"
	"errors"
	"time"
	"prom-lens-backend/internal/apps/authn/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Login(ctx context.Context, username, password string) (data model.User, err error)
	UpdateLastLogin(ctx context.Context, userID uint, at time.Time) error
	ChangeMyPassword(ctx context.Context, id uint, password string) error
	GetMe(ctx context.Context, id uint) (data model.User, err error)
	WithTx(tx *gorm.DB) AuthRepository
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Login(ctx context.Context, username, password string) (data model.User, err error) {
	const authErrorMsg = "账号或密码错误，请重新输入"
	if err = r.db.WithContext(ctx).Where("username = ?", username).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New(authErrorMsg)
		}
		return data, err
	}
	if data.IsActive != nil && !*data.IsActive {
		return data, errors.New(authErrorMsg)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(password)); err != nil {
		return data, errors.New(authErrorMsg)
	}
	return data, nil
}

func (r *authRepository) UpdateLastLogin(ctx context.Context, userID uint, at time.Time) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("last_login", at).Error
}

func (r *authRepository) ChangeMyPassword(ctx context.Context, id uint, password string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Select("password").Update("password", password).Error
}

func (r *authRepository) GetMe(ctx context.Context, id uint) (data model.User, err error) {
	if err = r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (r *authRepository) WithTx(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}
