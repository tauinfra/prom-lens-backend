package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceRepository interface {
	List(ctx context.Context, id uint) (namespaces []Namespace, err error)
	Get(ctx context.Context, id uint, name string) (namespace *v1.Namespace, err error)
	Create(ctx context.Context, id uint, body *v1.Namespace) (namespace *v1.Namespace, err error)
	Update(ctx context.Context, id uint, body *v1.Namespace) (namespace *v1.Namespace, err error)
	Delete(ctx context.Context, id uint, name string) error
}

type Namespace struct {
	Name      string            `json:"name"`
	Status    v1.NamespacePhase `json:"status"`
	Labels    map[string]string `json:"labels"`
	CreatedAt string            `json:"createdAt"`
}

type namespaces struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewNamespaceRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) NamespaceRepository {
	return &namespaces{cfgFactory: cfgFactory, gvkFactory: gvkFactory}
}

func (r *namespaces) List(ctx context.Context, id uint) (namespaces []Namespace, err error) {
	var ns Namespace
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 namespaces list failed. err: %v", err)
	}
	for _, item := range response.Items {
		ns.Name = item.Name
		ns.Status = item.Status.Phase
		ns.Labels = item.Labels
		ns.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
		namespaces = append(namespaces, ns)
	}
	r.gvkFactory.Complete(response)
	return namespaces, nil
}

func (r *namespaces) Create(ctx context.Context, id uint, body *v1.Namespace) (namespace *v1.Namespace, err error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Namespaces().Create(ctx, body, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 namespace create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *namespaces) Get(ctx context.Context, id uint, name string) (namespace *v1.Namespace, err error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 namespace get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *namespaces) Update(ctx context.Context, id uint, body *v1.Namespace) (namespace *v1.Namespace, err error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Namespaces().Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 namespace update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *namespaces) Delete(ctx context.Context, id uint, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	if err := client.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("kubernetes CoreV1 namespace delete failed. err: %v", err)
	}
	return nil
}
