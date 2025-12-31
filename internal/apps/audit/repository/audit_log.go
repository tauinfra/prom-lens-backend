package repository

import (
	"gorm.io/gorm"
	"valyria-backend/internal/apps/audit/model"
	pg "valyria-backend/internal/core/pagination"
)

// AuditLogRepository 定义接口
type AuditLogRepository interface {
	List(tx *gorm.DB, params pg.QueryParams) ([]model.AuditLog, pg.Pagination, error)
}

// auditLogRepository 实现了 AuditLogRepository 接口
type auditLogRepository struct {
	tx *gorm.DB
}

// NewAuditLogRepository 创建新的 AuditLogInterface 实例
func NewAuditLogRepository(tx *gorm.DB) AuditLogRepository {
	return &auditLogRepository{tx: tx}
}

// List 列表
func (r *auditLogRepository) List(tx *gorm.DB, params pg.QueryParams) (data []model.AuditLog, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(tx, &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}
