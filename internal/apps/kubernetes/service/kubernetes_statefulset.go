package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/repository"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
)

// StatefulSetService 定义接口
type StatefulSetService interface {
	List(ctx context.Context, id int, ns string) ([]repository.StatefulSet, error)
	Get(ctx context.Context, id int, ns, name string) (*appsv1.StatefulSet, error)
	GetDetail(ctx context.Context, id int, ns, name string) (repository.StatefulSet, error)
	Create(ctx context.Context, id int, ns string, body *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	Update(ctx context.Context, id int, ns string, body *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	Delete(ctx context.Context, id int, ns, name string) error
	Scale(ctx context.Context, id int, ns string, name string, replicas int32) (*autoscalingv1.Scale, error)
	Restart(ctx context.Context, id int, ns string, name string) (*appsv1.StatefulSet, error)
}

type statefulSetService struct {
	statefulSet repository.StatefulSetRepository
}

func NewStatefulSetService(statefulSet repository.StatefulSetRepository) StatefulSetService {
	return &statefulSetService{statefulSet: statefulSet}
}

// List 列表
func (s *statefulSetService) List(ctx context.Context, id int, ns string) ([]repository.StatefulSet, error) {
	return s.statefulSet.List(ctx, id, ns)
}

// Get 查询
func (s *statefulSetService) Get(ctx context.Context, id int, ns, name string) (*appsv1.StatefulSet, error) {
	return s.statefulSet.Get(ctx, id, ns, name)
}
func (s *statefulSetService) GetDetail(ctx context.Context, id int, ns, name string) (repository.StatefulSet, error) {
	return s.statefulSet.GetDetail(ctx, id, ns, name)
}

// Create 创建
func (s *statefulSetService) Create(ctx context.Context, id int, ns string, body *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return s.statefulSet.Create(ctx, id, ns, body)
}

// Update 更新
func (s *statefulSetService) Update(ctx context.Context, id int, ns string, body *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return s.statefulSet.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *statefulSetService) Delete(ctx context.Context, id int, ns, name string) (err error) {
	return s.statefulSet.Delete(ctx, id, ns, name)
}

// Scale 更新副本
func (s *statefulSetService) Scale(ctx context.Context, id int, ns string, name string, replicas int32) (*autoscalingv1.Scale, error) {
	return s.statefulSet.Scale(ctx, id, ns, name, replicas)
}

// Restart 重启副本
func (s *statefulSetService) Restart(ctx context.Context, id int, ns string, name string) (*appsv1.StatefulSet, error) {
	return s.statefulSet.Restart(ctx, id, ns, name)
}
