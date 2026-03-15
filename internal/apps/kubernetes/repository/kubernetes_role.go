package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RoleRepository interface {
	List(ctx context.Context, id uint, ns string) ([]Role, error)
	Get(ctx context.Context, id uint, ns, name string) (*rbacv1.Role, error)
	Create(ctx context.Context, id uint, ns string, body *rbacv1.Role) (*rbacv1.Role, error)
	Update(ctx context.Context, id uint, ns string, body *rbacv1.Role) (*rbacv1.Role, error)
	Delete(ctx context.Context, id uint, ns, name string) error
}

type Role struct {
	Name      string `json:"name"`
	Rules     int    `json:"rules"`
	CreatedAt string `json:"createdAt"`
}

type roleRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewRoleRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) RoleRepository {
	return &roleRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *roleRepository) List(ctx context.Context, id uint, ns string) (roles []Role, err error) {
	var role Role
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().Roles(ns).List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 roles list failed. err: %v", err)
	}
	for _, item := range response.Items {
		role.Name = item.Name
		role.Rules = len(item.Rules)
		role.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *roleRepository) Get(ctx context.Context, id uint, ns, name string) (*rbacv1.Role, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().Roles(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 role get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *roleRepository) Create(ctx context.Context, id uint, ns string, body *rbacv1.Role) (*rbacv1.Role, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().Roles(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 role create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *roleRepository) Update(ctx context.Context, id uint, ns string, body *rbacv1.Role) (*rbacv1.Role, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().Roles(ns).Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 role update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *roleRepository) Delete(ctx context.Context, id uint, ns, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.RbacV1().Roles(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes RbacV1 role delete failed. err: %v", err)
	}
	return nil
}
