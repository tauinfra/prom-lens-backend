package repository

import (
	"gorm.io/gorm"
	"valyria-backend/internal/apps/audit/model"
	pg "valyria-backend/internal/core/pagination"
)

// AuthLogRepository 定义接口
type AuthLogRepository interface {
	List(tx *gorm.DB, params pg.QueryParams) ([]model.AuthLog, pg.Pagination, error)
}

// authLogRepository 实现了 AuthLogRepository 接口
type authLogRepository struct {
	tx *gorm.DB
}

// NewAuthLogRepository 创建新的 AuthLogRepository 实例
func NewAuthLogRepository(tx *gorm.DB) AuthLogRepository {
	return &authLogRepository{tx: tx}
}

// List 列表
func (r *authLogRepository) List(tx *gorm.DB, params pg.QueryParams) (data []model.AuthLog, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(tx, &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}
