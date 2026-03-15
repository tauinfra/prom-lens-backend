package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TaskRunRepository interface {
	List(ctx context.Context, id uint, ns string) (taskRuns []TaskRun, err error)
	Get(ctx context.Context, id uint, ns, name string) (*v1.TaskRun, error)
	Create(ctx context.Context, id uint, ns string, body *v1.TaskRun) (*v1.TaskRun, error)
	Update(ctx context.Context, id uint, ns string, body *v1.TaskRun) (*v1.TaskRun, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}
type TaskRun struct {
	Name      string                 `json:"name,omitempty"`
	Namespace string                 `json:"namespace,omitempty"`
	Status    corev1.ConditionStatus `json:"status"`
	Reason    string                 `json:"reason,omitempty"`
	StartAt   string                 `json:"startAt,omitempty"`
	EndAt     string                 `json:"endAt,omitempty"`
	CreatedAt string                 `json:"createdAt,omitempty"`
}

type taskRunRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewTaskRunRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) TaskRunRepository {
	return &taskRunRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *taskRunRepository) List(ctx context.Context, id uint, ns string) (taskRuns []TaskRun, err error) {
	var taskRun TaskRun
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().TaskRuns(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return taskRuns, fmt.Errorf("kubernetes TektonV1 TaskRun list failed. err: %v", err)
	}
	for _, item := range response.Items {
		taskRun.Namespace = item.Namespace
		taskRun.Name = item.Name
		taskRun.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
		taskRuns = append(taskRuns, taskRun)
	}
	return taskRuns, nil
}

func (r *taskRunRepository) Get(ctx context.Context, id uint, ns, name string) (*v1.TaskRun, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().TaskRuns(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes TektonV1 TaskRun get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *taskRunRepository) Create(ctx context.Context, id uint, ns string, body *v1.TaskRun) (*v1.TaskRun, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().TaskRuns(ns).Create(ctx, body, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes TektonV1 TaskRun create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *taskRunRepository) Update(ctx context.Context, id uint, ns string, body *v1.TaskRun) (*v1.TaskRun, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().TaskRuns(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 TaskRun update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *taskRunRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.TektonV1().TaskRuns(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes TektonV1 TaskRun update failed. err: %v", err)
	}
	return nil
}
