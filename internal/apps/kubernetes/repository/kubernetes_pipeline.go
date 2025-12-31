package repository

import (
	"context"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/pkg/k8s/factory"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelineRepository interface {
	List(ctx context.Context, id int, ns string) (Pipelines []Pipeline, err error)
	Get(ctx context.Context, id int, ns, name string) (*v1.Pipeline, error)
	Create(ctx context.Context, id int, ns string, body *v1.Pipeline) (*v1.Pipeline, error)
	Update(ctx context.Context, id int, ns string, body *v1.Pipeline) (*v1.Pipeline, error)
	Delete(ctx context.Context, id int, ns, name string) error
}

type Pipeline struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

type pipelineRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewPipelineRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) PipelineRepository {
	return &pipelineRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *pipelineRepository) List(ctx context.Context, id int, ns string) (pipelines []Pipeline, err error) {
	var pipeline Pipeline
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().Pipelines(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Errorf("kubernetes TektonV1 Pipeline list failed. err: %v", err)
		return pipelines, err
	}
	for _, item := range response.Items {
		pipeline.Namespace = item.Namespace
		pipeline.Name = item.Name
		pipeline.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
		pipelines = append(pipelines, pipeline)
	}
	return pipelines, nil
}

func (r *pipelineRepository) Get(ctx context.Context, id int, ns, name string) (*v1.Pipeline, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().Pipelines(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("kubernetes TektonV1 Pipeline get failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *pipelineRepository) Create(ctx context.Context, id int, ns string, body *v1.Pipeline) (*v1.Pipeline, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().Pipelines(ns).Create(context.TODO(), body, metav1.CreateOptions{})
	if err != nil {
		logger.Errorf("kubernetes TektonV1 Pipeline create failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *pipelineRepository) Update(ctx context.Context, id int, ns string, body *v1.Pipeline) (*v1.Pipeline, error) {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.TektonV1().Pipelines(ns).Update(context.TODO(), body, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("kubernetes CoreV1 Pipeline update failed. err: %v", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *pipelineRepository) Delete(ctx context.Context, id int, ns, name string) error {
	client, err := r.cfgFactory.GetTektonClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.TektonV1().Pipelines(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logger.Errorf("kubernetes TektonV1 Pipeline update failed. err: %v", err)
	}
	return err
}
