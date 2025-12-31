package repository

import (
	"context"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/pkg/k8s/factory"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ReplicaSetRepository interface {
	List(ctx context.Context, id int, ns, labelSelector string) ([]ReplicaSet, error)
	Get(ctx context.Context, id int, ns, name string) (*appsv1.ReplicaSet, error)
	Update(ctx context.Context, id int, ns string, body *appsv1.ReplicaSet) (*appsv1.ReplicaSet, error)
	Delete(ctx context.Context, id int, ns, name string) error
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

func (r *replicaSetRepository) List(ctx context.Context, id int, ns, labelSelector string) (replicaSets []ReplicaSet, err error) {
	var replicaset ReplicaSet
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().ReplicaSets(ns).List(context.TODO(), metav1.ListOptions{
		LabelSelector:  labelSelector,
		TimeoutSeconds: &timeoutSeconds,
	})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 ReplicaSets list failed. err: %v", err)
		return nil, err
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

func (r *replicaSetRepository) Get(ctx context.Context, id int, ns, name string) (*appsv1.ReplicaSet, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().ReplicaSets(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 ReplicaSet get failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *replicaSetRepository) Update(ctx context.Context, id int, ns string, body *appsv1.ReplicaSet) (*appsv1.ReplicaSet, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.AppsV1().ReplicaSets(ns).Update(context.TODO(), body, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 ReplicaSet update failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *replicaSetRepository) Delete(ctx context.Context, id int, ns, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.AppsV1().ReplicaSets(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logger.Errorf("kubernetes AppsV1 ReplicaSet delete failed. err: %v", err)
		return err
	}
	return nil
}
