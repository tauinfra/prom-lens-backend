package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	corev1 "k8s.io/api/core/v1"
)

// PersistentVolumeClaimManager 定义接口
type PersistentVolumeClaimManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.PersistentVolumeClaim, error)
	Get(ctx context.Context, id uint, ns, name string) (*corev1.PersistentVolumeClaim, error)
	Create(ctx context.Context, id uint, ns string, body *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error)
	Update(ctx context.Context, id uint, ns string, body *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error)
	Delete(ctx context.Context, id uint, ns string, name string) error
}

type persistentVolumeClaimManager struct {
	persistentVolumeClaim repository.PersistentVolumeClaimRepository
}

func NewPersistentVolumeClaimManager(persistentVolumeClaim repository.PersistentVolumeClaimRepository) PersistentVolumeClaimManager {
	return &persistentVolumeClaimManager{persistentVolumeClaim: persistentVolumeClaim}
}

// List 列表
func (s *persistentVolumeClaimManager) List(ctx context.Context, id uint, ns string) ([]dto.PersistentVolumeClaim, error) {
	items, err := s.persistentVolumeClaim.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToPersistentVolumeClaimDTOs(items), nil
}

// Get 查询
func (s *persistentVolumeClaimManager) Get(ctx context.Context, id uint, ns, name string) (*corev1.PersistentVolumeClaim, error) {
	return s.persistentVolumeClaim.Get(ctx, id, ns, name)
}

// Create 创建
func (s *persistentVolumeClaimManager) Create(ctx context.Context, id uint, ns string, body *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	return s.persistentVolumeClaim.Create(ctx, id, ns, body)
}

// Update 更新
func (s *persistentVolumeClaimManager) Update(ctx context.Context, id uint, ns string, body *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	return s.persistentVolumeClaim.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *persistentVolumeClaimManager) Delete(ctx context.Context, id uint, ns, name string) (err error) {
	return s.persistentVolumeClaim.Delete(ctx, id, ns, name)
}
