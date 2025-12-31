package service

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"valyria-backend/internal/apps/kubernetes/repository"
)

// NodeService 定义接口
type NodeService interface {
	List(ctx context.Context, id int) ([]repository.Node, error)
	GetDetail(ctx context.Context, id int, name string) (*repository.Node, error)
	Get(ctx context.Context, id int, name string) (*corev1.Node, error)
	Update(ctx context.Context, id int, node *corev1.Node) (*corev1.Node, error)
	Cordon(ctx context.Context, id int, name string) (*corev1.Node, error)
}

type nodeService struct {
	node repository.NodeRepository
}

func NewNodeService(node repository.NodeRepository) NodeService {
	return &nodeService{node: node}
}

// List 列表
func (s *nodeService) List(ctx context.Context, id int) ([]repository.Node, error) {
	return s.node.List(ctx, id)
}

// Get 查询
func (s *nodeService) Get(ctx context.Context, id int, name string) (*corev1.Node, error) {
	return s.node.Get(ctx, id, name)
}

// GetDetail 查询
func (s *nodeService) GetDetail(ctx context.Context, id int, name string) (*repository.Node, error) {
	return s.node.GetDetail(ctx, id, name)
}

// Update 更新
func (s *nodeService) Update(ctx context.Context, id int, node *corev1.Node) (*corev1.Node, error) {
	return s.node.Update(ctx, id, node)
}

// Cordon 驱逐
func (s *nodeService) Cordon(ctx context.Context, id int, name string) (*corev1.Node, error) {
	return s.node.Cordon(ctx, id, name)
}
