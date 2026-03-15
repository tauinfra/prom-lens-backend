package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigmapRepository interface {
	List(ctx context.Context, id uint, ns string) ([]ConfigMap, error)
	Get(ctx context.Context, id uint, ns, name string) (*corev1.ConfigMap, error)
	Create(ctx context.Context, id uint, ns string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error)
	Update(ctx context.Context, id uint, ns string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type ConfigMap struct {
	Name      string `json:"name"`
	Data      int    `json:"data"`
	CreatedAt string `json:"createdAt"`
}

type configmapRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewConfigmapRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) ConfigmapRepository {
	return &configmapRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *configmapRepository) List(ctx context.Context, id uint, ns string) (configmaps []ConfigMap, err error) {
	var configmap ConfigMap
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().ConfigMaps(ns).List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 configmaps list failed. err: %v", err)
	}
	for _, item := range response.Items {
		configmap.Name = item.Name
		configmap.Data = len(item.Data) + len(item.BinaryData)                          // Key 数量
		configmap.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
		configmaps = append(configmaps, configmap)
	}
	return configmaps, nil
}

func (r *configmapRepository) Get(ctx context.Context, id uint, ns, name string) (*corev1.ConfigMap, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().ConfigMaps(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 configmap get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *configmapRepository) Create(ctx context.Context, id uint, ns string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().ConfigMaps(ns).Update(ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 configmaps create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *configmapRepository) Update(ctx context.Context, id uint, ns string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().ConfigMaps(ns).Update(ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 configmaps update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *configmapRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.CoreV1().ConfigMaps(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes CoreV1 configmap delete failed. err: %v", err)
	}
	return nil
}
