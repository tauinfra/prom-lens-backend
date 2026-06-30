package service

import (
	"context"
	"prom-lens-backend/internal/apps/audit/dto"
	"prom-lens-backend/internal/apps/audit/repository"
	pg "prom-lens-backend/internal/core/pagination"
)

type AuditLogService interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.AuditLogDTO, pg.Pagination, error)
}

// auditLogService 实现 AuditLogService 接口
type auditLogService struct {
	auditLog repository.AuditLogRepository
}

// NewAuditLogService 创建新的 AuditLogService 实例
func NewAuditLogService(auditLog repository.AuditLogRepository) AuditLogService {
	return &auditLogService{auditLog: auditLog}
}

// List 列表
func (s *auditLogService) List(ctx context.Context, params pg.QueryParams) ([]dto.AuditLogDTO, pg.Pagination, error) {
	logs, pagination, err := s.auditLog.List(ctx, params)
	if err != nil {
		return nil, pg.Pagination{}, err
	}
	data := make([]dto.AuditLogDTO, 0, len(logs))
	for _, log := range logs {
		data = append(data, dto.AuditLogDTO{
			ID:         log.ID,
			Username:   log.Username,
			UrlPath:    log.UrlPath,
			Method:     log.Method,
			IPAddress:  log.IPAddress,
			Agent:      log.Agent,
			StatusCode: log.StatusCode,
			Success:    log.Success,
			Params:     log.Params,
			Response:   log.Response,
			CreatedAt:  log.CreatedAt,
		})
	}
	return data, pagination, nil
}
