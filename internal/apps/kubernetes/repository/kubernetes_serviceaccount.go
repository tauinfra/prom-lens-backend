package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceAccountRepository interface {
	List(ctx context.Context, id uint, ns string) ([]ServiceAccount, error)
	Get(ctx context.Context, id uint, ns, name string) (*corev1.ServiceAccount, error)
	Create(ctx context.Context, id uint, ns string, body *corev1.ServiceAccount) (*corev1.ServiceAccount, error)
	Update(ctx context.Context, id uint, ns string, body *corev1.ServiceAccount) (*corev1.ServiceAccount, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type ServiceAccount struct {
	Name             string `json:"name"`
	Secrets          int    `json:"secrets"`
	ImagePullSecrets int    `json:"imagePullSecrets"`
	CreatedAt        string `json:"createdAt"`
}

type serviceAccountRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewServiceAccountRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) ServiceAccountRepository {
	return &serviceAccountRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *serviceAccountRepository) List(ctx context.Context, id uint, ns string) (accounts []ServiceAccount, err error) {
	var account ServiceAccount
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().ServiceAccounts(ns).List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 serviceaccounts list failed. err: %v", err)
	}
	for _, item := range response.Items {
		account.Name = item.Name
		account.Secrets = len(item.Secrets)
		account.ImagePullSecrets = len(item.ImagePullSecrets)
		account.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (r *serviceAccountRepository) Get(ctx context.Context, id uint, ns, name string) (*corev1.ServiceAccount, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().ServiceAccounts(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 serviceaccount get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *serviceAccountRepository) Create(ctx context.Context, id uint, ns string, body *corev1.ServiceAccount) (*corev1.ServiceAccount, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().ServiceAccounts(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 serviceaccount create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *serviceAccountRepository) Update(ctx context.Context, id uint, ns string, body *corev1.ServiceAccount) (*corev1.ServiceAccount, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().ServiceAccounts(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 serviceaccount update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *serviceAccountRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.CoreV1().ServiceAccounts(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes CoreV1 serviceaccount delete failed. err: %v", err)
	}
	return nil
}
