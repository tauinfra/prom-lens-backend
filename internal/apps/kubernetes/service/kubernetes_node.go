package service

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"
)

// NodeManager 定义接口
type NodeManager interface {
	List(ctx context.Context, id uint) ([]dto.Node, error)
	GetDetail(ctx context.Context, id uint, name string) (*dto.Node, error)
	Get(ctx context.Context, id uint, name string) (*corev1.Node, error)
	Update(ctx context.Context, id uint, node *corev1.Node) (*corev1.Node, error)
	Cordon(ctx context.Context, id uint, name string) (*corev1.Node, error)
}

type nodeManager struct {
	node repository.NodeRepository
}

func NewNodeManager(node repository.NodeRepository) NodeManager {
	return &nodeManager{node: node}
}

// List 列表
func (s *nodeManager) List(ctx context.Context, id uint) ([]dto.Node, error) {
	items, err := s.node.List(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.ToNodeDTOs(items), nil
}

// Get 查询
func (s *nodeManager) Get(ctx context.Context, id uint, name string) (*corev1.Node, error) {
	return s.node.Get(ctx, id, name)
}

// GetDetail 查询
func (s *nodeManager) GetDetail(ctx context.Context, id uint, name string) (*dto.Node, error) {
	item, err := s.node.GetDetail(ctx, id, name)
	if err != nil {
		return nil, err
	}
	dtoItem := dto.ToNodeDTO(*item)
	return &dtoItem, nil
}

// Update 更新
func (s *nodeManager) Update(ctx context.Context, id uint, node *corev1.Node) (*corev1.Node, error) {
	return s.node.Update(ctx, id, node)
}

// Cordon 驱逐
func (s *nodeManager) Cordon(ctx context.Context, id uint, name string) (*corev1.Node, error) {
	return s.node.Cordon(ctx, id, name)
}
