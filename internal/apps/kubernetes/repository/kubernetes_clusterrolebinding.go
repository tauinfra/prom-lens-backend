package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterRoleBindingRepository interface {
	List(ctx context.Context, id uint) ([]ClusterRoleBinding, error)
	Get(ctx context.Context, id uint, name string) (*rbacv1.ClusterRoleBinding, error)
	Create(ctx context.Context, id uint, body *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error)
	Update(ctx context.Context, id uint, body *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error)
	Delete(ctx context.Context, id uint, name string) error
}

type ClusterRoleBinding struct {
	Name      string `json:"name"`
	RoleRef   string `json:"roleRef"`
	Subjects  int    `json:"subjects"`
	CreatedAt string `json:"createdAt"`
}

type clusterRoleBindingRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewClusterRoleBindingRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) ClusterRoleBindingRepository {
	return &clusterRoleBindingRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *clusterRoleBindingRepository) List(ctx context.Context, id uint) (bindings []ClusterRoleBinding, err error) {
	var binding ClusterRoleBinding
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 clusterrolebindings list failed. err: %v", err)
	}
	for _, item := range response.Items {
		binding.Name = item.Name
		binding.RoleRef = item.RoleRef.Kind + "/" + item.RoleRef.Name
		binding.Subjects = len(item.Subjects)
		binding.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
		bindings = append(bindings, binding)
	}
	return bindings, nil
}

func (r *clusterRoleBindingRepository) Get(ctx context.Context, id uint, name string) (*rbacv1.ClusterRoleBinding, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().ClusterRoleBindings().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 clusterrolebinding get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *clusterRoleBindingRepository) Create(ctx context.Context, id uint, body *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().ClusterRoleBindings().Create(ctx, body, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 clusterrolebinding create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *clusterRoleBindingRepository) Update(ctx context.Context, id uint, body *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().ClusterRoleBindings().Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 clusterrolebinding update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *clusterRoleBindingRepository) Delete(ctx context.Context, id uint, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.RbacV1().ClusterRoleBindings().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes RbacV1 clusterrolebinding delete failed. err: %v", err)
	}
	return nil
}
