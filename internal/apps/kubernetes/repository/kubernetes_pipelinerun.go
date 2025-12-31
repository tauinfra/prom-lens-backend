package repository

import (
	"context"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/pkg/k8s/factory"
)

type PipelineRunRepository interface {
	List(ctx context.Context, id int, ns string) (PipelineRuns []PipelineRun, err error)
	Get(ctx context.Context, id int, ns, name string) (*v1.PipelineRun, error)
	Create(ctx context.Context, id int, ns string, body *v1.PipelineRun) (*v1.PipelineRun, error)
	Update(ctx context.Context, id int, ns string, body *v1.PipelineRun) (*v1.PipelineRun, error)
	Delete(ctx context.Context, id int, ns, name string) error
}

type PipelineRun struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

type pipelineRunRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewPipelineRunRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) PipelineRunRepository {
	return &pipelineRunRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *pipelineRunRepository) List(ctx context.Context, id int, ns string) (pipelineRuns []PipelineRun, err error) {
	var pipelineRun PipelineRun
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().PipelineRuns(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Errorf("kubernetes TektonV1 PipelineRun list failed. err: %v", err)
		return pipelineRuns, err
	}
	for _, item := range response.Items {
		pipelineRun.Namespace = item.Namespace
		pipelineRun.Name = item.Name
		pipelineRun.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
		pipelineRuns = append(pipelineRuns, pipelineRun)
	}
	return pipelineRuns, nil
}

func (r *pipelineRunRepository) Get(ctx context.Context, id int, ns, name string) (*v1.PipelineRun, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().PipelineRuns(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("kubernetes TektonV1 PipelineRun get failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *pipelineRunRepository) Create(ctx context.Context, id int, ns string, body *v1.PipelineRun) (*v1.PipelineRun, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().PipelineRuns(ns).Create(context.TODO(), body, metav1.CreateOptions{})
	if err != nil {
		logger.Errorf("kubernetes TektonV1 PipelineRun create failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *pipelineRunRepository) Update(ctx context.Context, id int, ns string, body *v1.PipelineRun) (*v1.PipelineRun, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().PipelineRuns(ns).Update(context.TODO(), body, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 PipelineRun update failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *pipelineRunRepository) Delete(ctx context.Context, id int, ns, name string) error {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.TektonV1().PipelineRuns(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logger.Errorf("kubernetes TektonV1 PipelineRun update failed. err: %v", err)
	}
	return err
}
