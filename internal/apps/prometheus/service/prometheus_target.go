package service

import (
	"context"
	"errors"
	"fmt"

	"prom-lens-backend/internal/apps/prometheus/dto"
	"prom-lens-backend/internal/apps/prometheus/executor"
	"prom-lens-backend/internal/apps/prometheus/model"
	"prom-lens-backend/internal/apps/prometheus/repository"
	"prom-lens-backend/internal/apps/prometheus/request"
	"prom-lens-backend/internal/core/logger"
	pg "prom-lens-backend/internal/core/pagination"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// TargetManager 定义服务层接口
type TargetManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.TargetDTO, pg.Pagination, error)
	Get(ctx context.Context, id int) (dto.TargetDTO, error)
	Create(ctx context.Context, groupID int, req *request.CreateTargetRequest) error
	Update(ctx context.Context, id int, req *request.UpdateTargetRequest) error
	Delete(ctx context.Context, id int) error
}

// targetManager 实现 TargetManager 接口
type targetManager struct {
	target repository.TargetRepository
	syncer executor.TargetSyncer
}

// NewTargetManager 创建新的 TargetManager 实例
func NewTargetManager(target repository.TargetRepository, syncer executor.TargetSyncer) TargetManager {
	return &targetManager{target: target, syncer: syncer}
}

// List 列表
func (s *targetManager) List(ctx context.Context, params pg.QueryParams) ([]dto.TargetDTO, pg.Pagination, error) {
	data, pagination, err := s.target.List(ctx, params)
	if err != nil {
		return nil, pg.Pagination{}, err
	}
	result := make([]dto.TargetDTO, 0, len(data))
	for _, item := range data {
		result = append(result, dto.ToTargetDTO(item))
	}
	return result, pagination, nil
}

func (s *targetManager) Get(ctx context.Context, id int) (dto.TargetDTO, error) {
	item, err := s.target.Get(ctx, id)
	if err != nil {
		return dto.TargetDTO{}, err
	}
	return dto.ToTargetDTO(item), nil
}

func (s *targetManager) Create(ctx context.Context, groupID int, req *request.CreateTargetRequest) error {
	data := &model.Target{
		GroupID:   groupID,
		IPAddress: req.IPAddress,
		Port:      req.Port,
		Labels:    req.Labels,
		Enabled:   req.Enabled,
	}
	if err := s.target.Create(ctx, data); err != nil {
		return translateTargetCreateError(err, groupID, req.IPAddress, req.Port)
	}
	if err := s.syncer.SyncTargetGroup(ctx, groupID); err != nil {
		if delErr := s.target.Delete(ctx, int(data.ID)); delErr != nil {
			logger.Errorf("[prom-sync] rollback target create id=%d failed err=%v", data.ID, delErr)
		}
		return fmt.Errorf("ConfigMap 同步失败，已回滚数据库写入: %w", err)
	}
	return nil
}

func (s *targetManager) Update(ctx context.Context, id int, req *request.UpdateTargetRequest) error {
	current, err := s.target.Get(ctx, id)
	if err != nil {
		return err
	}
	data := &model.Target{}
	if req.IPAddress != nil {
		data.IPAddress = *req.IPAddress
	}
	if req.Port != nil {
		data.Port = *req.Port
	}
	if req.Labels != nil {
		data.Labels = *req.Labels
	}
	if req.Enabled != nil {
		data.Enabled = req.Enabled
	}
	if err := s.target.Update(ctx, id, data); err != nil {
		return err
	}
	if err := s.syncer.SyncTargetGroup(ctx, current.GroupID); err != nil {
		if revErr := s.target.Update(ctx, id, targetSnapshot(&current)); revErr != nil {
			logger.Errorf("[prom-sync] rollback target update id=%d failed err=%v", id, revErr)
		}
		return fmt.Errorf("ConfigMap 同步失败，已回滚数据库更新: %w", err)
	}
	return nil
}

func (s *targetManager) Delete(ctx context.Context, id int) error {
	current, err := s.target.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.target.Delete(ctx, id); err != nil {
		return err
	}
	if err := s.syncer.SyncTargetGroup(ctx, current.GroupID); err != nil {
		restore := &model.Target{
			GroupID:   current.GroupID,
			IPAddress: current.IPAddress,
			Port:      current.Port,
			Labels:    current.Labels,
			Enabled:   current.Enabled,
		}
		if recErr := s.target.Create(ctx, restore); recErr != nil {
			logger.Errorf("[prom-sync] rollback target delete id=%d failed err=%v", id, recErr)
		}
		return fmt.Errorf("ConfigMap 同步失败，已回滚数据库删除: %w", err)
	}
	return nil
}

func targetSnapshot(t *model.Target) *model.Target {
	return &model.Target{
		IPAddress: t.IPAddress,
		Port:      t.Port,
		Labels:    t.Labels,
		Enabled:   t.Enabled,
	}
}

func translateTargetCreateError(err error, groupID int, ip string, port int) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return fmt.Errorf("采集节点 %s:%d 在组 %d 中已存在；若此前同步失败导致数据残留，请直接更新该节点或先删除后重建", ip, port, groupID)
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("采集节点 %s:%d 在组 %d 中已存在；若此前同步失败导致数据残留，请直接更新该节点或先删除后重建", ip, port, groupID)
	}
	return err
}
