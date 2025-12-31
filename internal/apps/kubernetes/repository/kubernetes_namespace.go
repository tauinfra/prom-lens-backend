package repository

import (
	"context"
	"valyria-backend/internal/pkg/k8s/factory"

	v1 "k8s.io/api/core/v1"
)

import (
	"fmt"
	"valyria-backend/internal/core/logger"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceRepository interface {
	List(ctx context.Context, id int) (namespaces []Namespace, err error)
	Create(ctx context.Context, id int, body *v1.Namespace) (namespace *v1.Namespace, err error)
}

type Namespace struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

type namespaces struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewNamespaceRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) NamespaceRepository {
	return &namespaces{cfgFactory: cfgFactory, gvkFactory: gvkFactory}
}

var timeoutSeconds int64 = 10

func (r *namespaces) List(ctx context.Context, id int) (namespaces []Namespace, err error) {
	var ns Namespace
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
		TimeoutSeconds: &timeoutSeconds,
	})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 namespaces list failed. err: %v", err)
		return nil, err
	}
	for _, item := range response.Items {
		ns.Name = item.Name
		ns.Status = fmt.Sprint(item.Status.Phase)
		ns.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
		namespaces = append(namespaces, ns)
	}
	r.gvkFactory.Complete(response)
	return namespaces, nil
}

func (r *namespaces) Create(ctx context.Context, id int, body *v1.Namespace) (namespace *v1.Namespace, err error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Namespaces().Create(context.TODO(), body, metav1.CreateOptions{})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 namespace create failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}
