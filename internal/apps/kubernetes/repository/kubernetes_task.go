package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TaskRepository interface {
	List(ctx context.Context, id uint, ns string) (Tasks []Task, err error)
	Get(ctx context.Context, id uint, ns, name string) (*v1.Task, error)
	Create(ctx context.Context, id uint, ns string, body *v1.Task) (*v1.Task, error)
	Update(ctx context.Context, id uint, ns string, body *v1.Task) (*v1.Task, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type Task struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

type taskRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewTaskRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) TaskRepository {
	return &taskRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *taskRepository) List(ctx context.Context, id uint, ns string) (tasks []Task, err error) {
	var task Task
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().Tasks(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return tasks, fmt.Errorf("kubernetes TektonV1 Task list failed. err: %v", err)
	}
	for _, item := range response.Items {
		task.Namespace = item.Namespace
		task.Name = item.Name
		task.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *taskRepository) Get(ctx context.Context, id uint, ns, name string) (*v1.Task, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().Tasks(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes TektonV1 Task get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *taskRepository) Create(ctx context.Context, id uint, ns string, body *v1.Task) (*v1.Task, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().Tasks(ns).Create(ctx, body, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes TektonV1 Task create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *taskRepository) Update(ctx context.Context, id uint, ns string, body *v1.Task) (*v1.Task, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().Tasks(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 Task update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *taskRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.TektonV1().Tasks(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes TektonV1 Task update failed. err: %v", err)
	}
	return nil
}
