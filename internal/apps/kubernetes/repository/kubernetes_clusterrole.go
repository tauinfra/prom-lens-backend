package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterRoleRepository interface {
	List(ctx context.Context, id uint) ([]ClusterRole, error)
	Get(ctx context.Context, id uint, name string) (*rbacv1.ClusterRole, error)
	Create(ctx context.Context, id uint, body *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error)
	Update(ctx context.Context, id uint, body *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error)
	Delete(ctx context.Context, id uint, name string) error
}

type ClusterRole struct {
	Name      string `json:"name"`
	Rules     int    `json:"rules"`
	CreatedAt string `json:"createdAt"`
}

type clusterRoleRepository struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewClusterRoleRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) ClusterRoleRepository {
	return &clusterRoleRepository{
		cfgFactory: cfgFactory,
		gvkFactory: gvkFactory,
	}
}

func (r *clusterRoleRepository) List(ctx context.Context, id uint) (roles []ClusterRole, err error) {
	var role ClusterRole
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{
		TimeoutSeconds: timeoutSeconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 clusterroles list failed. err: %v", err)
	}
	for _, item := range response.Items {
		role.Name = item.Name
		role.Rules = len(item.Rules)
		role.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *clusterRoleRepository) Get(ctx context.Context, id uint, name string) (*rbacv1.ClusterRole, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().ClusterRoles().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 clusterrole get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *clusterRoleRepository) Create(ctx context.Context, id uint, body *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().ClusterRoles().Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 clusterrole create failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *clusterRoleRepository) Update(ctx context.Context, id uint, body *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.RbacV1().ClusterRoles().Update(ctx, body, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes RbacV1 clusterrole update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *clusterRoleRepository) Delete(ctx context.Context, id uint, name string) error {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return err
	}
	err = client.RbacV1().ClusterRoles().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes RbacV1 clusterrole delete failed. err: %v", err)
	}
	return nil
}
