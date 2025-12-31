package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/repository"

	corev1 "k8s.io/api/core/v1"
)

// PersistentVolumeManager 定义接口
type PersistentVolumeManager interface {
	List(ctx context.Context, id int) ([]repository.PersistentVolume, error)
	Get(ctx context.Context, id int, name string) (*corev1.PersistentVolume, error)
	Create(ctx context.Context, id int, body *corev1.PersistentVolume) (*corev1.PersistentVolume, error)
	Update(ctx context.Context, id int, body *corev1.PersistentVolume) (*corev1.PersistentVolume, error)
	Delete(ctx context.Context, id int, name string) error
}
type persistentVolumeManager struct {
	persistentVolume repository.PersistentVolumeRepository
}

func NewPersistentVolumeManager(persistentVolume repository.PersistentVolumeRepository) PersistentVolumeManager {
	return &persistentVolumeManager{persistentVolume: persistentVolume}
}

// List 列表
func (s *persistentVolumeManager) List(ctx context.Context, id int) ([]repository.PersistentVolume, error) {
	return s.persistentVolume.List(ctx, id)
}

// Get 查询
func (s *persistentVolumeManager) Get(ctx context.Context, id int, name string) (*corev1.PersistentVolume, error) {
	return s.persistentVolume.Get(ctx, id, name)
}

// Create 创建
func (s *persistentVolumeManager) Create(ctx context.Context, id int, body *corev1.PersistentVolume) (*corev1.PersistentVolume, error) {
	return s.persistentVolume.Create(ctx, id, body)
}

// Update 更新
func (s *persistentVolumeManager) Update(ctx context.Context, id int, body *corev1.PersistentVolume) (*corev1.PersistentVolume, error) {
	return s.persistentVolume.Update(ctx, id, body)
}

// Delete 删除
func (s *persistentVolumeManager) Delete(ctx context.Context, id int, name string) (err error) {
	return s.persistentVolume.Delete(ctx, id, name)
}
