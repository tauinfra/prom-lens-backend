package service

import (
	"context"
	"errors"
	"strings"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/request"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
)

// StatefulSetManager 定义接口
type StatefulSetManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.StatefulSet, error)
	Get(ctx context.Context, id uint, ns, name string) (*appsv1.StatefulSet, error)
	GetDetail(ctx context.Context, id uint, ns, name string) (dto.StatefulSet, error)
	Create(ctx context.Context, id uint, ns string, body *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	Update(ctx context.Context, id uint, ns string, body *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	Delete(ctx context.Context, id uint, ns, name string) error
	DeleteBatch(ctx context.Context, id uint, ns string, req *request.BatchDeleteWorkloadRequest) error
	Scale(ctx context.Context, id uint, ns string, name string, replicas int32) (*autoscalingv1.Scale, error)
	Restart(ctx context.Context, id uint, ns string, name string) (*appsv1.StatefulSet, error)
}

type statefulSetManager struct {
	statefulSet repository.StatefulSetRepository
}

func NewStatefulSetManager(statefulSet repository.StatefulSetRepository) StatefulSetManager {
	return &statefulSetManager{statefulSet: statefulSet}
}

// List 列表
func (s *statefulSetManager) List(ctx context.Context, id uint, ns string) ([]dto.StatefulSet, error) {
	items, err := s.statefulSet.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToStatefulSetDTOs(items), nil
}

// Get 查询
func (s *statefulSetManager) Get(ctx context.Context, id uint, ns, name string) (*appsv1.StatefulSet, error) {
	return s.statefulSet.Get(ctx, id, ns, name)
}
func (s *statefulSetManager) GetDetail(ctx context.Context, id uint, ns, name string) (dto.StatefulSet, error) {
	item, err := s.statefulSet.GetDetail(ctx, id, ns, name)
	if err != nil {
		return dto.StatefulSet{}, err
	}
	return dto.ToStatefulSetDTO(item), nil
}

// Create 创建
func (s *statefulSetManager) Create(ctx context.Context, id uint, ns string, body *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return s.statefulSet.Create(ctx, id, ns, body)
}

// Update 更新
func (s *statefulSetManager) Update(ctx context.Context, id uint, ns string, body *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return s.statefulSet.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *statefulSetManager) Delete(ctx context.Context, id uint, ns, name string) (err error) {
	return s.statefulSet.Delete(ctx, id, ns, name)
}

func (s *statefulSetManager) DeleteBatch(ctx context.Context, id uint, ns string, req *request.BatchDeleteWorkloadRequest) error {
	if req == nil || len(req.Names) == 0 {
		return errors.New("statefulset names is required")
	}
	seen := make(map[string]struct{}, len(req.Names))
	for _, name := range req.Names {
		name = strings.TrimSpace(name)
		if name == "" {
			return errors.New("statefulset name is required")
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		if err := s.statefulSet.Delete(ctx, id, ns, name); err != nil {
			return err
		}
	}
	return nil
}

// Scale 更新副本
func (s *statefulSetManager) Scale(ctx context.Context, id uint, ns string, name string, replicas int32) (*autoscalingv1.Scale, error) {
	return s.statefulSet.Scale(ctx, id, ns, name, replicas)
}

// Restart 重启副本
func (s *statefulSetManager) Restart(ctx context.Context, id uint, ns string, name string) (*appsv1.StatefulSet, error) {
	return s.statefulSet.Restart(ctx, id, ns, name)
}
