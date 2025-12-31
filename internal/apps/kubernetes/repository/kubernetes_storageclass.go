package repository

import (
	"context"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/pkg/k8s/factory"
)

type StorageClassRepository interface {
	List(ctx context.Context, id int) ([]StorageClass, error)
	Get(ctx context.Context, id int, name string) (*storagev1.StorageClass, error)
	Create(ctx context.Context, id int, body *storagev1.StorageClass) (*storagev1.StorageClass, error)
	Update(ctx context.Context, id int, body *storagev1.StorageClass) (*storagev1.StorageClass, error)
	Delete(ctx context.Context, id int, name string) error
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

func (r *storageClassRepository) List(ctx context.Context, id int) (data []StorageClass, err error) {
	var storageClass StorageClass
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{
		TimeoutSeconds: &timeoutSeconds,
	})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 persistentVolumes list failed. err: %v", err)
		return nil, err
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

func (r *storageClassRepository) Get(ctx context.Context, id int, name string) (*storagev1.StorageClass, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.StorageV1().StorageClasses().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 persistentVolume get failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *storageClassRepository) Create(ctx context.Context, id int, body *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.StorageV1().StorageClasses().Update(context.TODO(), body, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 persistentVolume create failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *storageClassRepository) Update(ctx context.Context, id int, body *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.StorageV1().StorageClasses().Update(context.TODO(), body, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 persistentVolume update failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *storageClassRepository) Delete(ctx context.Context, id int, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.StorageV1().StorageClasses().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 persistentVolume delete failed. err: %v", err)
		return err
	}
	return nil
}
