package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/apps/dragon/adapter"

	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TektonTaskRepository interface {
	List(ctx context.Context, namespace string) ([]TektonTask, error)
	Get(ctx context.Context, namespace, name string) (*tektonv1.Task, error)
	Create(ctx context.Context, namespace string, body *tektonv1.Task) (*tektonv1.Task, error)
	Update(ctx context.Context, namespace string, body *tektonv1.Task) (*tektonv1.Task, error)
	Delete(ctx context.Context, namespace, name string) error
}

type TektonTask struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

type tektonTaskRepository struct{}

func NewTektonTaskRepository() TektonTaskRepository {
	return &tektonTaskRepository{}
}

func (r *tektonTaskRepository) List(ctx context.Context, namespace string) ([]TektonTask, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().Tasks(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton task list failed. err: %v", err)
	}
	items := make([]TektonTask, 0, len(response.Items))
	for _, item := range response.Items {
		items = append(items, TektonTask{
			Name:      item.Name,
			Namespace: item.Namespace,
			CreatedAt: item.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return items, nil
}

func (r *tektonTaskRepository) Get(ctx context.Context, namespace, name string) (*tektonv1.Task, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().Tasks(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton task get failed. err: %v", err)
	}
	return response, nil
}

func (r *tektonTaskRepository) Create(ctx context.Context, namespace string, body *tektonv1.Task) (*tektonv1.Task, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().Tasks(namespace).Create(ctx, body, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton task create failed. err: %v", err)
	}
	return response, nil
}

func (r *tektonTaskRepository) Update(ctx context.Context, namespace string, body *tektonv1.Task) (*tektonv1.Task, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().Tasks(namespace).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton task update failed. err: %v", err)
	}
	return response, nil
}

func (r *tektonTaskRepository) Delete(ctx context.Context, namespace, name string) error {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return err
	}
	if err := tektonClient.TektonV1().Tasks(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("tekton task delete failed. err: %v", err)
	}
	return nil
}
