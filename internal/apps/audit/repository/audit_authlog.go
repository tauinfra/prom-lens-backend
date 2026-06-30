package repository

import (
	"context"
	"prom-lens-backend/internal/apps/audit/model"
	pg "prom-lens-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// AuthLogRepository 定义接口
type AuthLogRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.AuthLog, pg.Pagination, error)
}

// authLogRepository 实现了 AuthLogRepository 接口
type authLogRepository struct {
	db *gorm.DB
}

// NewAuthLogRepository 创建新的 AuthLogRepository 实例
func NewAuthLogRepository(db *gorm.DB) AuthLogRepository {
	return &authLogRepository{db: db}
}

// List 列表
func (r *authLogRepository) List(ctx context.Context, params pg.QueryParams) (data []model.AuthLog, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}
