package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type JobRepository interface {
	List(ctx context.Context, id uint, ns string) ([]Job, error)
	Get(ctx context.Context, id uint, ns, name string) (*batchv1.Job, error)
	Create(ctx context.Context, id uint, ns string, body *batchv1.Job) (*batchv1.Job, error)
	Update(ctx context.Context, id uint, ns string, body *batchv1.Job) (*batchv1.Job, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type Job struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type jobRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewJobRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) JobRepository {
	return &jobRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *jobRepository) List(ctx context.Context, id uint, ns string) (jobs []Job, err error) {
	var job Job
	var client *kubernetes.Clientset
	client, err = r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.BatchV1().Jobs(ns).List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes BatchV1 jobs list failed. err: %v", err)
	}
	for _, item := range response.Items {
		job.Namespace = item.Namespace
		job.Name = item.Name
		job.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *jobRepository) Get(ctx context.Context, id uint, ns, name string) (job *batchv1.Job, err error) {
	var client *kubernetes.Clientset
	client, err = r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.BatchV1().Jobs(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes BatchV1 job get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *jobRepository) Create(ctx context.Context, id uint, ns string, body *batchv1.Job) (job *batchv1.Job, err error) {
	var client *kubernetes.Clientset
	client, err = r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.BatchV1().Jobs(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes BatchV1 job create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *jobRepository) Update(ctx context.Context, id uint, ns string, body *batchv1.Job) (job *batchv1.Job, err error) {
	var client *kubernetes.Clientset
	client, err = r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.BatchV1().Jobs(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes BatchV1 job update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *jobRepository) Delete(ctx context.Context, id uint, ns, name string) (err error) {
	var client *kubernetes.Clientset
	client, err = r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.BatchV1().Jobs(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes BatchV1 job delete failed. err: %v", err)
	}
	return nil
}
