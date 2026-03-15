package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PersistentVolumeRepository interface {
	List(ctx context.Context, id uint) ([]PersistentVolume, error)
	Get(ctx context.Context, id uint, name string) (*corev1.PersistentVolume, error)
	Create(ctx context.Context, id uint, body *corev1.PersistentVolume) (*corev1.PersistentVolume, error)
	Update(ctx context.Context, id uint, body *corev1.PersistentVolume) (*corev1.PersistentVolume, error)
	Delete(ctx context.Context, id uint, name string) error
}

type PersistentVolume struct {
	Name          string                               `json:"name"`
	Capacity      string                               `json:"capacity"`
	AccessModes   []corev1.PersistentVolumeAccessMode  `json:"accessModes"`
	ReclaimPolicy corev1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy"`
	Status        string                               `json:"status"`
	Claim         string                               `json:"claim"`
	StorageClass  string                               `json:"storageClass"`
	Reason        string                               `json:"reason"`
	CreatedAt     string                               `json:"createdAt"`
}

type persistentVolumeRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewPersistentVolumeRepository(kubeFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) PersistentVolumeRepository {
	return &persistentVolumeRepository{
		cfgFactory: kubeFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *persistentVolumeRepository) List(ctx context.Context, id uint) (data []PersistentVolume, err error) {
	var pv PersistentVolume
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 persistentVolumes list failed. err: %v", err)
	}
	for _, item := range response.Items {
		var claim string
		if item.Spec.ClaimRef != nil {
			claim = item.Spec.ClaimRef.Namespace + "/" + item.Spec.ClaimRef.Name
		}
		pv.Name = item.Name
		pv.Capacity = item.Spec.Capacity.Storage().String()
		pv.AccessModes = item.Spec.AccessModes
		pv.ReclaimPolicy = item.Spec.PersistentVolumeReclaimPolicy
		pv.Status = string(item.Status.Phase)
		pv.Claim = claim
		pv.StorageClass = item.Spec.StorageClassName
		pv.Reason = item.Status.Reason
		pv.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
		data = append(data, pv)
	}
	return data, nil
}

func (r *persistentVolumeRepository) Get(ctx context.Context, id uint, name string) (*corev1.PersistentVolume, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().PersistentVolumes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 persistentVolume get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *persistentVolumeRepository) Create(ctx context.Context, id uint, body *corev1.PersistentVolume) (*corev1.PersistentVolume, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().PersistentVolumes().Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 persistentVolume create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *persistentVolumeRepository) Update(ctx context.Context, id uint, body *corev1.PersistentVolume) (*corev1.PersistentVolume, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().PersistentVolumes().Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 persistentVolume update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *persistentVolumeRepository) Delete(ctx context.Context, id uint, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.CoreV1().PersistentVolumes().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes CoreV1 persistentVolume delete failed. err: %v", err)
	}
	return nil
}
