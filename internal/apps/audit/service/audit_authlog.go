package service

import (
	"gorm.io/gorm"
	"valyria-backend/internal/apps/audit/model"
	"valyria-backend/internal/apps/audit/repository"
	pg "valyria-backend/internal/core/pagination"
)

// AuthLogService 定义接口
type AuthLogService interface {
	List(params pg.QueryParams) ([]model.AuthLog, pg.Pagination, error)
}

// authLogService 实现 AuthLogInterface 接口
type authLogService struct {
	authLog repository.AuthLogRepository
	tx      *gorm.DB
}

// NewAuthLogService 创建新的 AuthLogService 实例
func NewAuthLogService(tx *gorm.DB, authLog repository.AuthLogRepository) AuthLogService {
	return &authLogService{tx: tx, authLog: authLog}
}

// List 列表
func (s *authLogService) List(params pg.QueryParams) ([]model.AuthLog, pg.Pagination, error) {
	return s.authLog.List(s.tx, params)
}
