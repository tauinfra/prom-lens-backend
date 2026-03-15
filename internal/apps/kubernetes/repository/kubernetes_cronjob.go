package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CronJobRepository interface {
	List(ctx context.Context, id uint, ns string) ([]CronJob, error)
	Get(ctx context.Context, id uint, ns, name string) (*batchv1.CronJob, error)
	Create(ctx context.Context, id uint, ns string, body *batchv1.CronJob) (*batchv1.CronJob, error)
	Update(ctx context.Context, id uint, ns string, body *batchv1.CronJob) (*batchv1.CronJob, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type CronJob struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type cronJobRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewCronJobRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) CronJobRepository {
	return &cronJobRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *cronJobRepository) List(ctx context.Context, id uint, ns string) (jobs []CronJob, err error) {
	var job CronJob
	var client *kubernetes.Clientset
	client, err = r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.BatchV1().CronJobs(ns).List(ctx, metav1.ListOptions{
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

func (r *cronJobRepository) Get(ctx context.Context, id uint, ns, name string) (job *batchv1.CronJob, err error) {
	var client *kubernetes.Clientset
	client, err = r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.BatchV1().CronJobs(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes BatchV1 job get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *cronJobRepository) Create(ctx context.Context, id uint, ns string, body *batchv1.CronJob) (cronJob *batchv1.CronJob, err error) {
	var client *kubernetes.Clientset
	client, err = r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.BatchV1().CronJobs(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes BatchV1 job create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *cronJobRepository) Update(ctx context.Context, id uint, ns string, body *batchv1.CronJob) (cronJob *batchv1.CronJob, err error) {
	var client *kubernetes.Clientset
	client, err = r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.BatchV1().CronJobs(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes BatchV1 job update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *cronJobRepository) Delete(ctx context.Context, id uint, ns, name string) (err error) {
	var client *kubernetes.Clientset
	client, err = r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.BatchV1().CronJobs(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes BatchV1 job delete failed. err: %v", err)
	}
	return nil
}
