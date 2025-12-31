package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/repository"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
)

// DeploymentService 定义接口
type DeploymentService interface {
	List(ctx context.Context, id int, ns string) ([]repository.Deployment, error)
	Get(ctx context.Context, id int, ns, name string) (*appsv1.Deployment, error)
	GetDetail(ctx context.Context, id int, ns, name string) (repository.Deployment, error)
	Create(ctx context.Context, id int, ns string, deployment *appsv1.Deployment) (*appsv1.Deployment, error)
	Update(ctx context.Context, id int, ns string, deployment *appsv1.Deployment) (*appsv1.Deployment, error)
	Delete(ctx context.Context, id int, ns string, name string) error
	Scale(ctx context.Context, id int, ns string, name string, replicas int32) (*autoscalingv1.Scale, error)
	Restart(ctx context.Context, id int, ns string, name string) (*appsv1.Deployment, error)
	Rollout(ctx context.Context, id int, ns string, name, rsName string) (*appsv1.Deployment, error)
}

type deploymentService struct {
	deployment repository.DeploymentRepository
}

func NewDeploymentService(deployment repository.DeploymentRepository) DeploymentService {
	return &deploymentService{deployment: deployment}
}

// List 列表
func (s *deploymentService) List(ctx context.Context, id int, ns string) ([]repository.Deployment, error) {
	return s.deployment.List(ctx, id, ns)
}

// Get 查询
func (s *deploymentService) Get(ctx context.Context, id int, ns, name string) (*appsv1.Deployment, error) {
	return s.deployment.Get(ctx, id, ns, name)
}
func (s *deploymentService) GetDetail(ctx context.Context, id int, ns, name string) (repository.Deployment, error) {
	return s.deployment.GetDetail(ctx, id, ns, name)
}

// Create 创建
func (s *deploymentService) Create(ctx context.Context, id int, ns string, body *appsv1.Deployment) (*appsv1.Deployment, error) {
	return s.deployment.Create(ctx, id, ns, body)
}

// Update 更新
func (s *deploymentService) Update(ctx context.Context, id int, ns string, body *appsv1.Deployment) (*appsv1.Deployment, error) {
	return s.deployment.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *deploymentService) Delete(ctx context.Context, id int, ns, name string) (err error) {
	return s.deployment.Delete(ctx, id, ns, name)
}

// Scale 更新副本
func (s *deploymentService) Scale(ctx context.Context, id int, ns string, name string, replicas int32) (*autoscalingv1.Scale, error) {
	return s.deployment.Scale(ctx, id, ns, name, replicas)
}

// Restart 重启副本
func (s *deploymentService) Restart(ctx context.Context, id int, ns string, name string) (*appsv1.Deployment, error) {
	return s.deployment.Restart(ctx, id, ns, name)
}

// Rollout 回滚副本
func (s *deploymentService) Rollout(ctx context.Context, id int, ns string, name, rsName string) (*appsv1.Deployment, error) {
	return s.deployment.Rollout(ctx, id, ns, name, rsName)
}
