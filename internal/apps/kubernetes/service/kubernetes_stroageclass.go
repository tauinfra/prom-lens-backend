package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/repository"

	storagev1 "k8s.io/api/storage/v1"
)

// StorageClassManager 定义接口
type StorageClassManager interface {
	List(ctx context.Context, id int) ([]repository.StorageClass, error)
	Get(ctx context.Context, id int, name string) (*storagev1.StorageClass, error)
	Create(ctx context.Context, id int, body *storagev1.StorageClass) (*storagev1.StorageClass, error)
	Update(ctx context.Context, id int, body *storagev1.StorageClass) (*storagev1.StorageClass, error)
	Delete(ctx context.Context, id int, name string) error
}
type storageClassManager struct {
	storageClass repository.StorageClassRepository
}

func NewStorageClassManager(storageClass repository.StorageClassRepository) StorageClassManager {
	return &storageClassManager{storageClass: storageClass}
}

// List 列表
func (s *storageClassManager) List(ctx context.Context, id int) ([]repository.StorageClass, error) {
	return s.storageClass.List(ctx, id)
}

// Get 查询
func (s *storageClassManager) Get(ctx context.Context, id int, name string) (*storagev1.StorageClass, error) {
	return s.storageClass.Get(ctx, id, name)
}

// Create 创建
func (s *storageClassManager) Create(ctx context.Context, id int, body *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	return s.storageClass.Create(ctx, id, body)
}

// Update 更新
func (s *storageClassManager) Update(ctx context.Context, id int, body *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	return s.storageClass.Update(ctx, id, body)
}

// Delete 删除
func (s *storageClassManager) Delete(ctx context.Context, id int, name string) (err error) {
	return s.storageClass.Delete(ctx, id, name)
}
