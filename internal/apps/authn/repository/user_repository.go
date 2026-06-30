package repository

import (
	"context"

	"prom-lens-backend/internal/apps/authn/model"
	pg "prom-lens-backend/internal/core/pagination"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.User, pg.Pagination, error)
	Get(ctx context.Context, id uint) (model.User, error)
	Create(ctx context.Context, data *model.User) error
	Update(ctx context.Context, id uint, data *model.User) error
	Delete(ctx context.Context, id uint) error
	ResetPassword(ctx context.Context, id uint, password string) error
	CountSuperusers(ctx context.Context) (int64, error)
	WithTx(tx *gorm.DB) UserRepository
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) List(ctx context.Context, params pg.QueryParams) ([]model.User, pg.Pagination, error) {
	var data []model.User
	pagination, err := pg.Paginate(r.db.WithContext(ctx), &data, params)
	if err != nil {
		return nil, pg.Pagination{}, err
	}
	return data, pagination, nil
}

func (r *userRepository) Get(ctx context.Context, id uint) (model.User, error) {
	var data model.User
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (r *userRepository) Create(ctx context.Context, data *model.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	data.Password = string(hash)
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *userRepository) Update(ctx context.Context, id uint, data *model.User) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(data).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *userRepository) ResetPassword(ctx context.Context, id uint, password string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("password", password).Error
}

func (r *userRepository) CountSuperusers(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("is_superuser = ?", true).Count(&count).Error
	return count, err
}

func (r *userRepository) WithTx(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
