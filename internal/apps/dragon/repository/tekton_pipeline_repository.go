package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/apps/dragon/adapter"

	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TektonPipelineRepository interface {
	List(ctx context.Context, namespace string) ([]TektonPipeline, error)
	Get(ctx context.Context, namespace, name string) (*tektonv1.Pipeline, error)
	Create(ctx context.Context, namespace string, body *tektonv1.Pipeline) (*tektonv1.Pipeline, error)
	Update(ctx context.Context, namespace string, body *tektonv1.Pipeline) (*tektonv1.Pipeline, error)
	Delete(ctx context.Context, namespace, name string) error
}

type TektonPipeline struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	CreatedAt string `json:"createdAt"`
}

type tektonPipelineRepository struct{}

func NewTektonPipelineRepository() TektonPipelineRepository {
	return &tektonPipelineRepository{}
}

func (r *tektonPipelineRepository) List(ctx context.Context, namespace string) ([]TektonPipeline, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().Pipelines(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton pipeline list failed. err: %v", err)
	}
	items := make([]TektonPipeline, 0, len(response.Items))
	for _, item := range response.Items {
		items = append(items, TektonPipeline{
			Name:      item.Name,
			Namespace: item.Namespace,
			CreatedAt: item.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return items, nil
}

func (r *tektonPipelineRepository) Get(ctx context.Context, namespace, name string) (*tektonv1.Pipeline, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().Pipelines(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton pipeline get failed. err: %v", err)
	}
	return response, nil
}

func (r *tektonPipelineRepository) Create(ctx context.Context, namespace string, body *tektonv1.Pipeline) (*tektonv1.Pipeline, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().Pipelines(namespace).Create(ctx, body, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton pipeline create failed. err: %v", err)
	}
	return response, nil
}

func (r *tektonPipelineRepository) Update(ctx context.Context, namespace string, body *tektonv1.Pipeline) (*tektonv1.Pipeline, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().Pipelines(namespace).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton pipeline update failed. err: %v", err)
	}
	return response, nil
}

func (r *tektonPipelineRepository) Delete(ctx context.Context, namespace, name string) error {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return err
	}
	if err := tektonClient.TektonV1().Pipelines(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("tekton pipeline delete failed. err: %v", err)
	}
	return nil
}
