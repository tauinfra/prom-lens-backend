package service

import (
	"context"
	"valyria-backend/internal/apps/audit/dto"
	"valyria-backend/internal/apps/audit/repository"
	pg "valyria-backend/internal/core/pagination"
)

// AuthLogService 定义接口
type AuthLogService interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.AuthLogDTO, pg.Pagination, error)
}

// authLogService 实现 AuthLogInterface 接口
type authLogService struct {
	authLog repository.AuthLogRepository
}

// NewAuthLogService 创建新的 AuthLogService 实例
func NewAuthLogService(authLog repository.AuthLogRepository) AuthLogService {
	return &authLogService{authLog: authLog}
}

// List 列表
func (s *authLogService) List(ctx context.Context, params pg.QueryParams) ([]dto.AuthLogDTO, pg.Pagination, error) {
	logs, pagination, err := s.authLog.List(ctx, params)
	if err != nil {
		return nil, pg.Pagination{}, err
	}
	data := make([]dto.AuthLogDTO, 0, len(logs))
	for _, log := range logs {
		data = append(data, dto.AuthLogDTO{
			ID:         log.ID,
			Username:   log.Username,
			IPAddress:  log.IPAddress,
			System:     log.System,
			Agent:      log.Agent,
			StatusCode: log.StatusCode,
			Success:    log.Success,
			ErrorCode:  log.ErrorCode,
			ErrorMsg:   log.ErrorMsg,
			CreatedAt:  log.CreatedAt,
		})
	}
	return data, pagination, nil
}
