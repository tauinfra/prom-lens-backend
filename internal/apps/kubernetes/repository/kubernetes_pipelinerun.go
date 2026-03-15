package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelineRunRepository interface {
	List(ctx context.Context, id uint, ns string) (PipelineRuns []PipelineRun, err error)
	Get(ctx context.Context, id uint, ns, name string) (*v1.PipelineRun, error)
	Create(ctx context.Context, id uint, ns string, body *v1.PipelineRun) (*v1.PipelineRun, error)
	Update(ctx context.Context, id uint, ns string, body *v1.PipelineRun) (*v1.PipelineRun, error)
	Delete(ctx context.Context, id uint, ns, name string) error
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

func (r *pipelineRunRepository) List(ctx context.Context, id uint, ns string) (pipelineRuns []PipelineRun, err error) {
	var pipelineRun PipelineRun
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().PipelineRuns(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return pipelineRuns, fmt.Errorf("kubernetes TektonV1 PipelineRun list failed. err: %v", err)
	}
	for _, item := range response.Items {
		pipelineRun.Namespace = item.Namespace
		pipelineRun.Name = item.Name
		pipelineRun.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
		pipelineRuns = append(pipelineRuns, pipelineRun)
	}
	return pipelineRuns, nil
}

func (r *pipelineRunRepository) Get(ctx context.Context, id uint, ns, name string) (*v1.PipelineRun, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().PipelineRuns(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes TektonV1 PipelineRun get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *pipelineRunRepository) Create(ctx context.Context, id uint, ns string, body *v1.PipelineRun) (*v1.PipelineRun, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().PipelineRuns(ns).Create(ctx, body, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes TektonV1 PipelineRun create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *pipelineRunRepository) Update(ctx context.Context, id uint, ns string, body *v1.PipelineRun) (*v1.PipelineRun, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().PipelineRuns(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 PipelineRun update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *pipelineRunRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.TektonV1().PipelineRuns(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes TektonV1 PipelineRun update failed. err: %v", err)
	}
	return nil
}
