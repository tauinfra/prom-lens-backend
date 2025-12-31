package service

import (
	"gorm.io/gorm"
	"valyria-backend/internal/apps/audit/model"
	"valyria-backend/internal/apps/audit/repository"
	pg "valyria-backend/internal/core/pagination"
)

type AuditLogService interface {
	List(params pg.QueryParams) ([]model.AuditLog, pg.Pagination, error)
}

// auditLogService 实现 AuditLogService 接口
type auditLogService struct {
	auditLog repository.AuditLogRepository
	tx       *gorm.DB
}

// NewAuditLogService 创建新的 AuditLogService 实例
func NewAuditLogService(tx *gorm.DB, auditLog repository.AuditLogRepository) AuditLogService {
	return &auditLogService{tx: tx, auditLog: auditLog}
}

// List 列表
func (s *auditLogService) List(params pg.QueryParams) ([]model.AuditLog, pg.Pagination, error) {
	return s.auditLog.List(s.tx, params)
}
