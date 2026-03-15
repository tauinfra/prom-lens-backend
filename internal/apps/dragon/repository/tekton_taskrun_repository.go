package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/apps/dragon/adapter"

	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TektonTaskRunRepository interface {
	List(ctx context.Context, namespace string) ([]TektonTaskRun, error)
	Get(ctx context.Context, namespace, name string) (*tektonv1.TaskRun, error)
	Create(ctx context.Context, namespace string, body *tektonv1.TaskRun) (*tektonv1.TaskRun, error)
	Update(ctx context.Context, namespace string, body *tektonv1.TaskRun) (*tektonv1.TaskRun, error)
	Delete(ctx context.Context, namespace, name string) error
	BatchDelete(ctx context.Context, namespace string, names []string) error
}

type TektonTaskRun struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Succeeded string `json:"succeeded,omitempty"`
	Reason    string `json:"reason,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

type tektonTaskRunRepository struct{}

func NewTektonTaskRunRepository() TektonTaskRunRepository {
	return &tektonTaskRunRepository{}
}

func (r *tektonTaskRunRepository) List(ctx context.Context, namespace string) ([]TektonTaskRun, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().TaskRuns(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton taskrun list failed. err: %v", err)
	}
	items := make([]TektonTaskRun, 0, len(response.Items))
	for _, item := range response.Items {
		succeeded, reason := tektonSucceededAndReason(item.Status.Conditions)
		items = append(items, TektonTaskRun{
			Name:      item.Name,
			Namespace: item.Namespace,
			Succeeded: succeeded,
			Reason:    reason,
			CreatedAt: item.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return items, nil
}

func (r *tektonTaskRunRepository) Get(ctx context.Context, namespace, name string) (*tektonv1.TaskRun, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().TaskRuns(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton taskrun get failed. err: %v", err)
	}
	return response, nil
}

func (r *tektonTaskRunRepository) Create(ctx context.Context, namespace string, body *tektonv1.TaskRun) (*tektonv1.TaskRun, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().TaskRuns(namespace).Create(ctx, body, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton taskrun create failed. err: %v", err)
	}
	return response, nil
}

func (r *tektonTaskRunRepository) Update(ctx context.Context, namespace string, body *tektonv1.TaskRun) (*tektonv1.TaskRun, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().TaskRuns(namespace).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton taskrun update failed. err: %v", err)
	}
	return response, nil
}

func (r *tektonTaskRunRepository) Delete(ctx context.Context, namespace, name string) error {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return err
	}
	if err := tektonClient.TektonV1().TaskRuns(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("tekton taskrun delete failed. err: %v", err)
	}
	return nil
}

func (r *tektonTaskRunRepository) BatchDelete(ctx context.Context, namespace string, names []string) error {
	if len(names) == 0 {
		return nil
	}
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return err
	}
	for _, name := range names {
		if name == "" {
			continue
		}
		if err := tektonClient.TektonV1().TaskRuns(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
			return fmt.Errorf("tekton taskrun delete failed. name=%s err: %v", name, err)
		}
	}
	return nil
}
