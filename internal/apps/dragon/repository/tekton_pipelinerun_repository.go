package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/apps/dragon/adapter"

	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TektonPipelineRunRepository interface {
	List(ctx context.Context, namespace string) ([]TektonPipelineRun, error)
	Get(ctx context.Context, namespace, name string) (*tektonv1.PipelineRun, error)
	Create(ctx context.Context, namespace string, body *tektonv1.PipelineRun) (*tektonv1.PipelineRun, error)
	Update(ctx context.Context, namespace string, body *tektonv1.PipelineRun) (*tektonv1.PipelineRun, error)
	Delete(ctx context.Context, namespace, name string) error
	BatchDelete(ctx context.Context, namespace string, names []string) error
}

type TektonPipelineRun struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Succeeded string `json:"succeeded,omitempty"`
	Reason    string `json:"reason,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

type tektonPipelineRunRepository struct{}

func NewTektonPipelineRunRepository() TektonPipelineRunRepository {
	return &tektonPipelineRunRepository{}
}

func (r *tektonPipelineRunRepository) List(ctx context.Context, namespace string) ([]TektonPipelineRun, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().PipelineRuns(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton pipelinerun list failed. err: %v", err)
	}
	items := make([]TektonPipelineRun, 0, len(response.Items))
	for _, item := range response.Items {
		succeeded, reason := tektonSucceededAndReason(item.Status.Conditions)
		items = append(items, TektonPipelineRun{
			Name:      item.Name,
			Namespace: item.Namespace,
			Succeeded: succeeded,
			Reason:    reason,
			CreatedAt: item.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return items, nil
}

func (r *tektonPipelineRunRepository) Get(ctx context.Context, namespace, name string) (*tektonv1.PipelineRun, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().PipelineRuns(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton pipelinerun get failed. err: %v", err)
	}
	return response, nil
}

func (r *tektonPipelineRunRepository) Create(ctx context.Context, namespace string, body *tektonv1.PipelineRun) (*tektonv1.PipelineRun, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().PipelineRuns(namespace).Create(ctx, body, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton pipelinerun create failed. err: %v", err)
	}
	return response, nil
}

func (r *tektonPipelineRunRepository) Update(ctx context.Context, namespace string, body *tektonv1.PipelineRun) (*tektonv1.PipelineRun, error) {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	response, err := tektonClient.TektonV1().PipelineRuns(namespace).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("tekton pipelinerun update failed. err: %v", err)
	}
	return response, nil
}

func (r *tektonPipelineRunRepository) Delete(ctx context.Context, namespace, name string) error {
	tektonClient, _, err := adapter.NewTektonClients()
	if err != nil {
		return err
	}
	if err := tektonClient.TektonV1().PipelineRuns(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("tekton pipelinerun delete failed. err: %v", err)
	}
	return nil
}

func (r *tektonPipelineRunRepository) BatchDelete(ctx context.Context, namespace string, names []string) error {
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
		if err := tektonClient.TektonV1().PipelineRuns(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
			return fmt.Errorf("tekton pipelinerun delete failed. name=%s err: %v", name, err)
		}
	}
	return nil
}
