package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IngressRepository interface {
	List(ctx context.Context, id uint, ns string) ([]Ingress, error)
	Get(ctx context.Context, id uint, ns, name string) (*networkingv1.Ingress, error)
	Create(ctx context.Context, id uint, ns string, body *networkingv1.Ingress) (*networkingv1.Ingress, error)
	Update(ctx context.Context, id uint, ns string, body *networkingv1.Ingress) (*networkingv1.Ingress, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type Ingress struct {
	Name      string   `json:"name"`
	Hosts     []string `json:"hosts"`
	Paths     int      `json:"paths"`
	TLS       bool     `json:"tls"`
	Endpoint  string   `json:"endpoint"`
	ClassName string   `json:"className"`
	CreatedAt string   `json:"createdAt"`
}

type ingressRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewIngressRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) IngressRepository {
	return &ingressRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *ingressRepository) List(ctx context.Context, id uint, ns string) (ingresses []Ingress, err error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.NetworkingV1().Ingresses(ns).List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes NetworkingV1 ingresses list failed. err: %v", err)
	}

	for _, item := range response.Items {
		hosts := make([]string, 0, len(item.Spec.Rules))
		seenHosts := make(map[string]struct{}, len(item.Spec.Rules))
		pathsCount := 0
		for _, rule := range item.Spec.Rules {
			host := rule.Host
			if host == "" {
				host = "*"
			}
			if rule.HTTP != nil {
				pathsCount += len(rule.HTTP.Paths)
			}
			if _, exists := seenHosts[host]; !exists {
				seenHosts[host] = struct{}{}
				hosts = append(hosts, host)
			}
		}
		ingress := Ingress{
			Name:      item.Name,
			Hosts:     hosts,
			Paths:     pathsCount,
			TLS:       len(item.Spec.TLS) > 0,
			CreatedAt: item.CreationTimestamp.UTC().Format("2006-01-02T15:04:05Z"),
		}
		for _, lb := range item.Status.LoadBalancer.Ingress {
			if lb.Hostname != "" {
				ingress.Endpoint = lb.Hostname
				break
			}
		}
		if item.Spec.IngressClassName != nil {
			ingress.ClassName = *item.Spec.IngressClassName
		}
		ingresses = append(ingresses, ingress)
	}
	return ingresses, nil
}

func (r *ingressRepository) Get(ctx context.Context, id uint, ns, name string) (*networkingv1.Ingress, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.NetworkingV1().Ingresses(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes NetworkingV1 ingress get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *ingressRepository) Create(ctx context.Context, id uint, ns string, body *networkingv1.Ingress) (*networkingv1.Ingress, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.NetworkingV1().Ingresses(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes NetworkingV1 ingress create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *ingressRepository) Update(ctx context.Context, id uint, ns string, body *networkingv1.Ingress) (*networkingv1.Ingress, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.NetworkingV1().Ingresses(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes NetworkingV1 ingress update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *ingressRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.NetworkingV1().Ingresses(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes NetworkingV1 ingress delete failed. err: %v", err)
	}
	return nil
}
