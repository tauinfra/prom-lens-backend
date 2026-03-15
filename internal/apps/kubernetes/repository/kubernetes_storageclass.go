package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StorageClassRepository interface {
	List(ctx context.Context, id uint) ([]StorageClass, error)
	Get(ctx context.Context, id uint, name string) (*storagev1.StorageClass, error)
	Create(ctx context.Context, id uint, body *storagev1.StorageClass) (*storagev1.StorageClass, error)
	Update(ctx context.Context, id uint, body *storagev1.StorageClass) (*storagev1.StorageClass, error)
	Delete(ctx context.Context, id uint, name string) error
}

type StorageClass struct {
	Name                 string `json:"name"`
	Provisioner          string `json:"provisioner"`
	ReclaimPolicy        string `json:"reclaimPolicy"`
	VolumeBindingMode    string `json:"volumeBindingMode"`
	AllowVolumeExpansion bool   `json:"allowVolumeExpansion"`
	CreatedAt            string `json:"createdAt"`
}

type storageClassRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewStorageClassRepository(kubeFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) StorageClassRepository {
	return &storageClassRepository{
		cfgFactory: kubeFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *storageClassRepository) List(ctx context.Context, id uint) (data []StorageClass, err error) {
	var storageClass StorageClass
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 persistentVolumes list failed. err: %v", err)
	}
	for _, item := range response.Items {
		storageClass.Name = item.Name
		storageClass.Provisioner = item.Provisioner
		storageClass.ReclaimPolicy = string(*item.ReclaimPolicy)
		storageClass.VolumeBindingMode = string(*item.VolumeBindingMode)
		storageClass.AllowVolumeExpansion = item.AllowVolumeExpansion != nil && *item.AllowVolumeExpansion
		storageClass.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
		data = append(data, storageClass)
	}
	return data, nil
}

func (r *storageClassRepository) Get(ctx context.Context, id uint, name string) (*storagev1.StorageClass, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.StorageV1().StorageClasses().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 persistentVolume get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *storageClassRepository) Create(ctx context.Context, id uint, body *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.StorageV1().StorageClasses().Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 persistentVolume create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *storageClassRepository) Update(ctx context.Context, id uint, body *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.StorageV1().StorageClasses().Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 persistentVolume update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *storageClassRepository) Delete(ctx context.Context, id uint, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.StorageV1().StorageClasses().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes CoreV1 persistentVolume delete failed. err: %v", err)
	}
	return nil
}
