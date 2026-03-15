package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PersistentVolumeClaimRepository interface {
	List(ctx context.Context, id uint, ns string) ([]PersistentVolumeClaim, error)
	Get(ctx context.Context, id uint, ns, name string) (*corev1.PersistentVolumeClaim, error)
	Create(ctx context.Context, id uint, ns string, body *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error)
	Update(ctx context.Context, id uint, ns string, body *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type PersistentVolumeClaim struct {
	Name         string                              `json:"name"`
	Capacity     string                              `json:"capacity"`
	AccessModes  []corev1.PersistentVolumeAccessMode `json:"accessModes"`
	Status       string                              `json:"status"`
	StorageClass *string                             `json:"storageClass"`
	Volume       string                              `json:"volume"`
	CreatedAt    string                              `json:"createdAt"`
}

type persistentVolumeClaimRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewPersistentVolumeClaimRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) PersistentVolumeClaimRepository {
	return &persistentVolumeClaimRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *persistentVolumeClaimRepository) List(ctx context.Context, id uint, ns string) (data []PersistentVolumeClaim, err error) {
	var pvc PersistentVolumeClaim
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().PersistentVolumeClaims(ns).List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 persistentVolumeClaims list failed. err: %v", err)
	}
	for _, item := range response.Items {
		pvc.Name = item.Name
		pvc.Volume = item.Spec.VolumeName
		pvc.AccessModes = item.Status.AccessModes
		pvc.Capacity = item.Status.Capacity.Storage().String()
		pvc.Status = string(item.Status.Phase)
		pvc.StorageClass = item.Spec.StorageClassName
		pvc.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
		data = append(data, pvc)
	}
	return data, nil
}

func (r *persistentVolumeClaimRepository) Get(ctx context.Context, id uint, ns, name string) (*corev1.PersistentVolumeClaim, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().PersistentVolumeClaims(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 persistentVolumeClaim get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *persistentVolumeClaimRepository) Create(ctx context.Context, id uint, ns string, body *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().PersistentVolumeClaims(ns).Create(ctx, body, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 persistentVolumeClaim create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *persistentVolumeClaimRepository) Update(ctx context.Context, id uint, ns string, body *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().PersistentVolumeClaims(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 persistentVolumeClaim update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *persistentVolumeClaimRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.CoreV1().PersistentVolumeClaims(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes CoreV1 persistentVolumeClaim delete failed. err: %v", err)
	}
	return nil
}
