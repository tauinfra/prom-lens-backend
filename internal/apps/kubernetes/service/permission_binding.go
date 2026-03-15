package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceAll 表示集群级权限，对应 ClusterRoleBinding
const NamespaceAll = "*"

// Binding 名称格式：ClusterRoleBinding = val-{username}-{roleName}，RoleBinding = val-{username}-{namespace}-{roleName}
const bindingNamePrefix = "val-"

var (
	reK8sNameInvalid = regexp.MustCompile(`[^a-z0-9-]`)
	reK8sNameHyphen  = regexp.MustCompile(`-+`)
)

// sanitizeForK8sName 将字符串转为合法 K8s 资源名（小写、连字符）
func sanitizeForK8sName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, ":", "-")
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, ".", "-")
	s = reK8sNameInvalid.ReplaceAllString(s, "-")
	s = reK8sNameHyphen.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "default"
	}
	return s
}

// roleNameForBinding 用于 Binding 名称的 role 部分：移除 val: 前缀，只保留 admin/developer 等
// 如 val:admin -> admin，val:developer -> developer
func roleNameForBinding(role string) string {
	if strings.HasPrefix(role, "val:") {
		return strings.TrimPrefix(role, "val:")
	}
	return role
}

// ClusterRoleBindingName 生成 ClusterRoleBinding 名称：val-{username}-{roleName}（roleName 已去掉 val: 前缀）
func ClusterRoleBindingName(username, role string) string {
	return bindingNamePrefix + sanitizeForK8sName(username) + "-" + sanitizeForK8sName(roleNameForBinding(role))
}

// RoleBindingName 生成 RoleBinding 名称：val-{username}-{namespace}-{roleName}（roleName 已去掉 val: 前缀）
func RoleBindingName(username, namespace, role string) string {
	return bindingNamePrefix + sanitizeForK8sName(username) + "-" + sanitizeForK8sName(namespace) + "-" + sanitizeForK8sName(roleNameForBinding(role))
}

// IsClusterScoped 是否集群级：namespace == "*" 使用 ClusterRoleBinding，否则使用 RoleBinding
func IsClusterScoped(namespace string) bool {
	return namespace == NamespaceAll
}

// BuildClusterRoleBindingForPermission 根据权限与用户名构建 ClusterRoleBinding
// 示例：name=val-username-admin, subjects User name=username, roleRef ClusterRole name=role
func BuildClusterRoleBindingForPermission(name, username, role string) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "User",
				Name: username,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     role,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}

// BuildRoleBindingForPermission 根据权限与用户名构建 RoleBinding
// 示例：name=val-username-default-developer, namespace=default, subjects User name=username, roleRef ClusterRole name=role
func BuildRoleBindingForPermission(name, namespace, username, role string) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "User",
				Name: username,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     role,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}

// SyncClusterRoleBindingCreate 在集群中创建 ClusterRoleBinding（namespace == "*" 时使用）
func SyncClusterRoleBindingCreate(ctx context.Context, clusterID uint, username, role string, createFn func(ctx context.Context, id uint, body *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error)) error {
	name := ClusterRoleBindingName(username, role)
	body := BuildClusterRoleBindingForPermission(name, username, role)
	_, err := createFn(ctx, clusterID, body)
	if err != nil {
		return fmt.Errorf("create clusterrolebinding %s: %w", name, err)
	}
	return nil
}

// SyncClusterRoleBindingDelete 在集群中删除 ClusterRoleBinding
func SyncClusterRoleBindingDelete(ctx context.Context, clusterID uint, username, role string, deleteFn func(ctx context.Context, id uint, name string) error) error {
	name := ClusterRoleBindingName(username, role)
	return deleteFn(ctx, clusterID, name)
}

// SyncRoleBindingCreate 在集群指定命名空间中创建 RoleBinding（namespace != "*" 时使用）
func SyncRoleBindingCreate(ctx context.Context, clusterID uint, namespace string, username, role string, createFn func(ctx context.Context, id uint, ns string, body *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error)) error {
	name := RoleBindingName(username, namespace, role)
	body := BuildRoleBindingForPermission(name, namespace, username, role)
	_, err := createFn(ctx, clusterID, namespace, body)
	if err != nil {
		return fmt.Errorf("create rolebinding %s/%s: %w", namespace, name, err)
	}
	return nil
}

// SyncRoleBindingDelete 在集群指定命名空间中删除 RoleBinding
func SyncRoleBindingDelete(ctx context.Context, clusterID uint, namespace string, username, role string, deleteFn func(ctx context.Context, id uint, ns, name string) error) error {
	name := RoleBindingName(username, namespace, role)
	return deleteFn(ctx, clusterID, namespace, name)
}
