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

// DeploymentManager 定义接口
type DeploymentManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.Deployment, error)
	Get(ctx context.Context, id uint, ns, name string) (*appsv1.Deployment, error)
	GetDetail(ctx context.Context, id uint, ns, name string) (dto.Deployment, error)
	Create(ctx context.Context, id uint, ns string, deployment *appsv1.Deployment) (*appsv1.Deployment, error)
	Update(ctx context.Context, id uint, ns string, deployment *appsv1.Deployment) (*appsv1.Deployment, error)
	Delete(ctx context.Context, id uint, ns string, name string) error
	DeleteBatch(ctx context.Context, id uint, ns string, req *request.BatchDeleteWorkloadRequest) error
	Scale(ctx context.Context, id uint, ns string, name string, replicas int32) (*autoscalingv1.Scale, error)
	Restart(ctx context.Context, id uint, ns string, name string) (*appsv1.Deployment, error)
	Rollout(ctx context.Context, id uint, ns string, name, rsName string) (*appsv1.Deployment, error)
}

type deploymentManager struct {
	deployment repository.DeploymentRepository
}

func NewDeploymentManager(deployment repository.DeploymentRepository) DeploymentManager {
	return &deploymentManager{deployment: deployment}
}

// List 列表
func (s *deploymentManager) List(ctx context.Context, id uint, ns string) ([]dto.Deployment, error) {
	items, err := s.deployment.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToDeploymentDTOs(items), nil
}

// Get 查询
func (s *deploymentManager) Get(ctx context.Context, id uint, ns, name string) (*appsv1.Deployment, error) {
	return s.deployment.Get(ctx, id, ns, name)
}
func (s *deploymentManager) GetDetail(ctx context.Context, id uint, ns, name string) (dto.Deployment, error) {
	item, err := s.deployment.GetDetail(ctx, id, ns, name)
	if err != nil {
		return dto.Deployment{}, err
	}
	return dto.ToDeploymentDTO(item), nil
}

// Create 创建
func (s *deploymentManager) Create(ctx context.Context, id uint, ns string, body *appsv1.Deployment) (*appsv1.Deployment, error) {
	return s.deployment.Create(ctx, id, ns, body)
}

// Update 更新
func (s *deploymentManager) Update(ctx context.Context, id uint, ns string, body *appsv1.Deployment) (*appsv1.Deployment, error) {
	return s.deployment.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *deploymentManager) Delete(ctx context.Context, id uint, ns, name string) (err error) {
	return s.deployment.Delete(ctx, id, ns, name)
}

func (s *deploymentManager) DeleteBatch(ctx context.Context, id uint, ns string, req *request.BatchDeleteWorkloadRequest) error {
	if req == nil || len(req.Names) == 0 {
		return errors.New("deployment names is required")
	}
	seen := make(map[string]struct{}, len(req.Names))
	for _, name := range req.Names {
		name = strings.TrimSpace(name)
		if name == "" {
			return errors.New("deployment name is required")
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		if err := s.deployment.Delete(ctx, id, ns, name); err != nil {
			return err
		}
	}
	return nil
}

// Scale 更新副本
func (s *deploymentManager) Scale(ctx context.Context, id uint, ns string, name string, replicas int32) (*autoscalingv1.Scale, error) {
	return s.deployment.Scale(ctx, id, ns, name, replicas)
}

// Restart 重启副本
func (s *deploymentManager) Restart(ctx context.Context, id uint, ns string, name string) (*appsv1.Deployment, error) {
	return s.deployment.Restart(ctx, id, ns, name)
}

// Rollout 回滚副本
func (s *deploymentManager) Rollout(ctx context.Context, id uint, ns string, name, rsName string) (*appsv1.Deployment, error) {
	return s.deployment.Rollout(ctx, id, ns, name, rsName)
}
