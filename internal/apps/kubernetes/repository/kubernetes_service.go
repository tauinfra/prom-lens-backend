package repository

import (
	"context"
	"fmt"
	"strings"
	"valyria-backend/internal/pkg/k8s/factory"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceRepository interface {
	List(ctx context.Context, id uint, ns, labelSelector string) (services []Service, err error)
	Get(ctx context.Context, id uint, ns, name string) (*corev1.Service, error)
	Create(ctx context.Context, id uint, ns string, body *corev1.Service) (*corev1.Service, error)
	Update(ctx context.Context, id uint, ns string, body *corev1.Service) (*corev1.Service, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type Service struct {
	Name      string            `json:"name,omitempty"`
	Type      corev1.ServiceType `json:"type,omitempty"`
	ClusterIP string            `json:"clusterIP,omitempty"`
	Ports     string            `json:"ports,omitempty"`
	CreatedAt string            `json:"createdAt,omitempty"`
}

type serviceRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewServiceRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) ServiceRepository {
	return &serviceRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *serviceRepository) List(ctx context.Context, id uint, ns, labelSelector string) (services []Service, err error) {
	var service Service
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Services(ns).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return services, fmt.Errorf("kubernetes CoreV1 service list failed. err: %v", err)
	}
	for _, item := range response.Items {
		var ports string
		for _, port := range item.Spec.Ports {
			if port.NodePort != 0 {
				ports += fmt.Sprintf("%d:%d/%v,", port.Port, port.NodePort, port.Protocol)
			} else {
				ports += fmt.Sprintf("%d/%v,", port.Port, port.Protocol)
			}
		}
		service = Service{
			Name:      item.Name,
			Type:      item.Spec.Type,
			ClusterIP: item.Spec.ClusterIP,
			Ports:     strings.TrimRight(ports, ","),
			CreatedAt: item.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
		}
		services = append(services, service)
	}
	return services, nil
}

func (r *serviceRepository) Get(ctx context.Context, id uint, ns, name string) (*corev1.Service, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Services(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 service get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *serviceRepository) Create(ctx context.Context, id uint, ns string, body *corev1.Service) (*corev1.Service, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Services(ns).Create(ctx, body, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 service create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *serviceRepository) Update(ctx context.Context, id uint, ns string, body *corev1.Service) (*corev1.Service, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Services(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 service update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *serviceRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.CoreV1().Services(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes CoreV1 service update failed. err: %v", err)
	}
	return nil
}
