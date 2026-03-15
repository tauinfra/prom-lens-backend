package repository

import (
	"context"
	"fmt"
	"log"
	"valyria-backend/internal/pkg/k8s/factory"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SecretRepository interface {
	List(ctx context.Context, id uint, ns string) ([]Secret, error)
	Get(ctx context.Context, id uint, ns, name string) (*corev1.Secret, error)
	Create(ctx context.Context, id uint, ns string, body *corev1.Secret) (*corev1.Secret, error)
	Update(ctx context.Context, id uint, ns string, body *corev1.Secret) (*corev1.Secret, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type Secret struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Data      int    `json:"data"`
	CreatedAt string `json:"createdAt"`
}

type secretRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewSecretRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) SecretRepository {
	return &secretRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *secretRepository) List(ctx context.Context, id uint, ns string) (secrets []Secret, err error) {
	var secret Secret
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Secrets(ns).List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 secrets list failed. err: %v", err)
	}
	for _, item := range response.Items {
		secret.Name = item.Name
		secret.Type = string(item.Type)
		secret.Data = len(item.Data)
		secret.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05") // 格式化时间
		secrets = append(secrets, secret)
	}
	return secrets, nil
}

func (r *secretRepository) Get(ctx context.Context, id uint, ns, name string) (*corev1.Secret, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Secrets(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 secret get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *secretRepository) Create(ctx context.Context, id uint, ns string, body *corev1.Secret) (*corev1.Secret, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Secrets(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		log.Println("kubernetes CoreV1 secret create failed. err:", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *secretRepository) Update(ctx context.Context, id uint, ns string, body *corev1.Secret) (*corev1.Secret, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Secrets(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		log.Println("kubernetes CoreV1 secret update failed. err:", err)
		return nil, err
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *secretRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.CoreV1().ConfigMaps(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes CoreV1 secret delete failed. err: %v", err)
	}
	return nil
}
