package service

import (
	"context"
	"errors"
	"strings"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/request"

	appsv1 "k8s.io/api/apps/v1"
)

// DaemonSetManager 定义接口
type DaemonSetManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.DaemonSet, error)
	Get(ctx context.Context, id uint, ns, name string) (*appsv1.DaemonSet, error)
	GetDetail(ctx context.Context, id uint, ns, name string) (dto.DaemonSet, error)
	Create(ctx context.Context, id uint, ns string, body *appsv1.DaemonSet) (*appsv1.DaemonSet, error)
	Update(ctx context.Context, id uint, ns string, body *appsv1.DaemonSet) (*appsv1.DaemonSet, error)
	Delete(ctx context.Context, id uint, ns string, name string) error
	DeleteBatch(ctx context.Context, id uint, ns string, req *request.BatchDeleteWorkloadRequest) error
}

type daemonSetManager struct {
	daemonset repository.DaemonSetRepository
}

func NewDaemonSetManager(daemonset repository.DaemonSetRepository) DaemonSetManager {
	return &daemonSetManager{daemonset: daemonset}
}

// List 列表
func (s *daemonSetManager) List(ctx context.Context, id uint, ns string) ([]dto.DaemonSet, error) {
	items, err := s.daemonset.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToDaemonSetDTOs(items), nil
}

// Get 查询
func (s *daemonSetManager) Get(ctx context.Context, id uint, ns, name string) (*appsv1.DaemonSet, error) {
	return s.daemonset.Get(ctx, id, ns, name)
}

// GetDetail 详情
func (s *daemonSetManager) GetDetail(ctx context.Context, id uint, ns, name string) (dto.DaemonSet, error) {
	item, err := s.daemonset.GetDetail(ctx, id, ns, name)
	if err != nil {
		return dto.DaemonSet{}, err
	}
	return dto.ToDaemonSetDTO(item), nil
}

// Create 创建
func (s *daemonSetManager) Create(ctx context.Context, id uint, ns string, body *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return s.daemonset.Create(ctx, id, ns, body)
}

// Update 更新
func (s *daemonSetManager) Update(ctx context.Context, id uint, ns string, body *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return s.daemonset.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *daemonSetManager) Delete(ctx context.Context, id uint, ns, name string) (err error) {
	return s.daemonset.Delete(ctx, id, ns, name)
}

func (s *daemonSetManager) DeleteBatch(ctx context.Context, id uint, ns string, req *request.BatchDeleteWorkloadRequest) error {
	if req == nil || len(req.Names) == 0 {
		return errors.New("daemonset names is required")
	}
	seen := make(map[string]struct{}, len(req.Names))
	for _, name := range req.Names {
		name = strings.TrimSpace(name)
		if name == "" {
			return errors.New("daemonset name is required")
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		if err := s.daemonset.Delete(ctx, id, ns, name); err != nil {
			return err
		}
	}
	return nil
}
