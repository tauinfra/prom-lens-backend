package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ReplicaSetRepository interface {
	List(ctx context.Context, id uint, ns, labelSelector string) ([]ReplicaSet, error)
	Get(ctx context.Context, id uint, ns, name string) (*appsv1.ReplicaSet, error)
	Update(ctx context.Context, id uint, ns string, body *appsv1.ReplicaSet) (*appsv1.ReplicaSet, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type ReplicaSet struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	Revision  string `json:"revision"`
	CreatedAt string `json:"createdAt"`
}

type replicaSetRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewReplicaSetRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) ReplicaSetRepository {
	return &replicaSetRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *replicaSetRepository) List(ctx context.Context, id uint, ns, labelSelector string) (replicaSets []ReplicaSet, err error) {
	var replicaset ReplicaSet
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().ReplicaSets(ns).List(ctx, metav1.ListOptions{
		LabelSelector:  labelSelector,
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes AppsV1 ReplicaSets list failed. err: %v", err)
	}
	for _, item := range response.Items {
		replicaset.Namespace = item.Namespace
		replicaset.Name = item.Name
		replicaset.Image = item.Spec.Template.Spec.Containers[0].Image                   // 默认获取第一个容器镜像版本
		replicaset.Revision = item.Annotations["deployment.kubernetes.io/revision"]      // 修订版本(默认 10 个版本)
		replicaset.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
		replicaSets = append(replicaSets, replicaset)
	}
	return
}

func (r *replicaSetRepository) Get(ctx context.Context, id uint, ns, name string) (*appsv1.ReplicaSet, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().ReplicaSets(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes AppsV1 ReplicaSet get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *replicaSetRepository) Update(ctx context.Context, id uint, ns string, body *appsv1.ReplicaSet) (*appsv1.ReplicaSet, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().ReplicaSets(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes AppsV1 ReplicaSet update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *replicaSetRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.AppsV1().ReplicaSets(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes AppsV1 ReplicaSet delete failed. err: %v", err)
	}
	return nil
}
