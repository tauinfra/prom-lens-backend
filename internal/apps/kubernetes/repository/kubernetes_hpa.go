package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HpaRepository interface {
	List(ctx context.Context, id uint, ns string) ([]HPA, error)
	Get(ctx context.Context, id uint, ns, name string) (*autoscalingv2.HorizontalPodAutoscaler, error)
	Create(ctx context.Context, id uint, ns string, hpa *autoscalingv2.HorizontalPodAutoscaler) (*autoscalingv2.HorizontalPodAutoscaler, error)
	Update(ctx context.Context, id uint, ns string, hpa *autoscalingv2.HorizontalPodAutoscaler) (*autoscalingv2.HorizontalPodAutoscaler, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type HPA struct {
	Namespace       string `json:"namespace"`
	Name            string `json:"name"`
	ScaleTargetKind string `json:"scaleTargetKind"`
	ScaleTargetName string `json:"scaleTargetName"`
	MinReplicas     *int32 `json:"minReplicas,omitempty"`
	MaxReplicas     int32  `json:"maxReplicas"`
	CurrentReplicas int32  `json:"currentReplicas"`
	DesiredReplicas int32  `json:"desiredReplicas"`
	CreatedAt       string `json:"createdAt"`
}

type hpaRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewHpaRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) HpaRepository {
	return &hpaRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *hpaRepository) List(ctx context.Context, id uint, ns string) (list []HPA, err error) {
	client, err := r.cfgFactory.GetClientSetAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return nil, err
	}
	response, err := client.AutoscalingV2().HorizontalPodAutoscalers(ns).List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes AutoscalingV2 horizontalpodautoscalers list failed. err: %v", err)
	}
	for _, item := range response.Items {
		list = append(list, r.toHPA(item))
	}
	return list, nil
}

func (r *hpaRepository) Get(ctx context.Context, id uint, ns, name string) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	client, err := r.cfgFactory.GetClientSetAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return nil, err
	}
	obj, err := client.AutoscalingV2().HorizontalPodAutoscalers(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes AutoscalingV2 horizontalpodautoscaler get failed. err: %v", err)
	}
	r.gvkFactory.Complete(obj)
	return obj, nil
}

func (r *hpaRepository) Create(ctx context.Context, id uint, ns string, hpa *autoscalingv2.HorizontalPodAutoscaler) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	client, err := r.cfgFactory.GetClientSetAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return nil, err
	}
	if hpa.Namespace == "" {
		hpa.Namespace = ns
	}
	obj, err := client.AutoscalingV2().HorizontalPodAutoscalers(ns).Create(ctx, hpa, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes AutoscalingV2 horizontalpodautoscaler create failed. err: %v", err)
	}
	r.gvkFactory.Complete(obj)
	return obj, nil
}

func (r *hpaRepository) Update(ctx context.Context, id uint, ns string, hpa *autoscalingv2.HorizontalPodAutoscaler) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	client, err := r.cfgFactory.GetClientSetAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return nil, err
	}
	obj, err := client.AutoscalingV2().HorizontalPodAutoscalers(ns).Update(ctx, hpa, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes AutoscalingV2 horizontalpodautoscaler update failed. err: %v", err)
	}
	r.gvkFactory.Complete(obj)
	return obj, nil
}

func (r *hpaRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetClientSetAsUser(ctx, id, ImpersonateUsername(ctx))
	if err != nil {
		return err
	}
	err = client.AutoscalingV2().HorizontalPodAutoscalers(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes AutoscalingV2 horizontalpodautoscaler delete failed. err: %v", err)
	}
	return nil
}

func (r *hpaRepository) toHPA(item autoscalingv2.HorizontalPodAutoscaler) HPA {
	out := HPA{
		Namespace:        item.Namespace,
		Name:              item.Name,
		ScaleTargetKind:   item.Spec.ScaleTargetRef.Kind,
		ScaleTargetName:   item.Spec.ScaleTargetRef.Name,
		MaxReplicas:       item.Spec.MaxReplicas,
		CurrentReplicas:   item.Status.CurrentReplicas,
		DesiredReplicas:   item.Status.DesiredReplicas,
		CreatedAt:         item.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
	}
	if item.Spec.MinReplicas != nil {
		out.MinReplicas = item.Spec.MinReplicas
	}
	return out
}
