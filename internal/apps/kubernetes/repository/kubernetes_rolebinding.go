package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RoleBindingRepository interface {
	List(ctx context.Context, id uint, ns string) ([]RoleBinding, error)
	Get(ctx context.Context, id uint, ns, name string) (*rbacv1.RoleBinding, error)
	Create(ctx context.Context, id uint, ns string, body *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error)
	Update(ctx context.Context, id uint, ns string, body *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type RoleBinding struct {
	Name      string `json:"name"`
	RoleRef   string `json:"roleRef"`
	Subjects  int    `json:"subjects"`
	CreatedAt string `json:"createdAt"`
}

type roleBindingRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewRoleBindingRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) RoleBindingRepository {
	return &roleBindingRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *roleBindingRepository) List(ctx context.Context, id uint, ns string) (bindings []RoleBinding, err error) {
	var binding RoleBinding
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().RoleBindings(ns).List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 rolebindings list failed. err: %v", err)
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

func (r *roleBindingRepository) Get(ctx context.Context, id uint, ns, name string) (*rbacv1.RoleBinding, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().RoleBindings(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 rolebinding get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *roleBindingRepository) Create(ctx context.Context, id uint, ns string, body *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().RoleBindings(ns).Create(ctx, body, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 rolebinding create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *roleBindingRepository) Update(ctx context.Context, id uint, ns string, body *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().RoleBindings(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 rolebinding update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *roleBindingRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.RbacV1().RoleBindings(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes RbacV1 rolebinding delete failed. err: %v", err)
	}
	return nil
}
