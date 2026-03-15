package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IngressClassRepository interface {
	List(ctx context.Context, id uint) ([]IngressClass, error)
	Get(ctx context.Context, id uint, name string) (*networkingv1.IngressClass, error)
	Create(ctx context.Context, id uint, body *networkingv1.IngressClass) (*networkingv1.IngressClass, error)
	Update(ctx context.Context, id uint, body *networkingv1.IngressClass) (*networkingv1.IngressClass, error)
	Delete(ctx context.Context, id uint, name string) error
}

type IngressClass struct {
	Name       string `json:"name"`
	Controller string `json:"controller"`
	Parameters string `json:"parameters"`
	CreatedAt  string `json:"createdAt"`
}

type ingressClassRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewIngressClassRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) IngressClassRepository {
	return &ingressClassRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *ingressClassRepository) List(ctx context.Context, id uint) (classes []IngressClass, err error) {
	var class IngressClass
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.NetworkingV1().IngressClasses().List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes NetworkingV1 ingressclasses list failed. err: %v", err)
	}
	for _, item := range response.Items {
		class.Name = item.Name
		class.Controller = item.Spec.Controller
		if item.Spec.Parameters != nil {
			class.Parameters = item.Spec.Parameters.Kind + "/" + item.Spec.Parameters.Name
		} else {
			class.Parameters = ""
		}
		class.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
		classes = append(classes, class)
	}
	return classes, nil
}

func (r *ingressClassRepository) Get(ctx context.Context, id uint, name string) (*networkingv1.IngressClass, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.NetworkingV1().IngressClasses().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes NetworkingV1 ingressclass get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *ingressClassRepository) Create(ctx context.Context, id uint, body *networkingv1.IngressClass) (*networkingv1.IngressClass, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.NetworkingV1().IngressClasses().Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes NetworkingV1 ingressclass create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *ingressClassRepository) Update(ctx context.Context, id uint, body *networkingv1.IngressClass) (*networkingv1.IngressClass, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.NetworkingV1().IngressClasses().Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes NetworkingV1 ingressclass update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *ingressClassRepository) Delete(ctx context.Context, id uint, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.NetworkingV1().IngressClasses().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes NetworkingV1 ingressclass delete failed. err: %v", err)
	}
	return nil
}
