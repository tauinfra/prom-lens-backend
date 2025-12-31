package repository

import (
	"context"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/pkg/k8s/factory"
)

type TaskRepository interface {
	List(ctx context.Context, id int, ns string) (Tasks []Task, err error)
	Get(ctx context.Context, id int, ns, name string) (*v1.Task, error)
	Create(ctx context.Context, id int, ns string, body *v1.Task) (*v1.Task, error)
	Update(ctx context.Context, id int, ns string, body *v1.Task) (*v1.Task, error)
	Delete(ctx context.Context, id int, ns, name string) error
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

func (r *taskRepository) List(ctx context.Context, id int, ns string) (tasks []Task, err error) {
	var task Task
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().Tasks(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Errorf("kubernetes TektonV1 Task list failed. err: %v", err)
		return tasks, err
	}
	for _, item := range response.Items {
		task.Namespace = item.Namespace
		task.Name = item.Name
		task.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *taskRepository) Get(ctx context.Context, id int, ns, name string) (*v1.Task, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().Tasks(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("kubernetes TektonV1 Task get failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *taskRepository) Create(ctx context.Context, id int, ns string, body *v1.Task) (*v1.Task, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().Tasks(ns).Create(context.TODO(), body, metav1.CreateOptions{})
	if err != nil {
		logger.Errorf("kubernetes TektonV1 Task create failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *taskRepository) Update(ctx context.Context, id int, ns string, body *v1.Task) (*v1.Task, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().Tasks(ns).Update(context.TODO(), body, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 Task update failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *taskRepository) Delete(ctx context.Context, id int, ns, name string) error {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.TektonV1().Tasks(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logger.Errorf("kubernetes TektonV1 Task update failed. err: %v", err)
	}
	return err
}
