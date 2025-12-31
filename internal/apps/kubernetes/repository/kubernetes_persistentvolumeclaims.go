package repository

import (
	"context"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/pkg/k8s/factory"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PersistentVolumeClaimRepository interface {
	List(ctx context.Context, id int, ns string) ([]PersistentVolumeClaim, error)
	Get(ctx context.Context, id int, ns, name string) (*corev1.PersistentVolumeClaim, error)
	Create(ctx context.Context, id int, ns string, body *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error)
	Update(ctx context.Context, id int, ns string, body *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error)
	Delete(ctx context.Context, id int, ns, name string) error
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

func (r *persistentVolumeClaimRepository) List(ctx context.Context, id int, ns string) (data []PersistentVolumeClaim, err error) {
	var pvc PersistentVolumeClaim
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().PersistentVolumeClaims(ns).List(context.TODO(), metav1.ListOptions{
		TimeoutSeconds: &timeoutSeconds,
	})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 persistentVolumeClaims list failed. err: %v", err)
		return nil, err
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

func (r *persistentVolumeClaimRepository) Get(ctx context.Context, id int, ns, name string) (*corev1.PersistentVolumeClaim, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().PersistentVolumeClaims(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 persistentVolumeClaim get failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *persistentVolumeClaimRepository) Create(ctx context.Context, id int, ns string, body *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().PersistentVolumeClaims(ns).Update(context.TODO(), body, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 persistentVolumeClaim create failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *persistentVolumeClaimRepository) Update(ctx context.Context, id int, ns string, body *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().PersistentVolumeClaims(ns).Update(context.TODO(), body, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 persistentVolumeClaim update failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *persistentVolumeClaimRepository) Delete(ctx context.Context, id int, ns, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.CoreV1().PersistentVolumeClaims(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 persistentVolumeClaim delete failed. err: %v", err)
		return err
	}
	return nil
}
