package repository

import (
	"context"
	"prom-lens-backend/internal/apps/audit/model"
	pg "prom-lens-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// AuditLogRepository 定义接口
type AuditLogRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.AuditLog, pg.Pagination, error)
}

// auditLogRepository 实现了 AuditLogRepository 接口
type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository 创建新的 AuditLogInterface 实例
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

// List 列表
func (r *auditLogRepository) List(ctx context.Context, params pg.QueryParams) (data []model.AuditLog, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}
