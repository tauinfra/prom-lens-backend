package routes

import (
	"valyria-backend/internal/pkg/di"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func KubernetesRouters(rg *gin.RouterGroup, tx *gorm.DB, provider *di.Provider) {
	kube := rg.Group("/kubernetes")
	{
		// Permissions
		kube.GET("/permissions",
			provider.PermissionManager.RequirePermission("kubernetes:permission:list"),
			provider.Kubernetes.Permission.Controller.List)
		kube.POST("/permissions",
			provider.PermissionManager.RequirePermission("kubernetes:permission:create"),
			provider.Kubernetes.Permission.Controller.Create)
		kube.POST("/permissions/batch",
			provider.PermissionManager.RequirePermission("kubernetes:permission:create"),
			provider.Kubernetes.Permission.Controller.CreateBatch)
		kube.DELETE("/permissions/batch-delete",
			provider.PermissionManager.RequirePermission("kubernetes:permission:delete"),
			provider.Kubernetes.Permission.Controller.DeleteBatch)
		kube.DELETE("/permissions/:id",
			provider.PermissionManager.RequirePermission("kubernetes:permission:delete"),
			provider.Kubernetes.Permission.Controller.Delete)
		// Cluster Routes
		kube.GET("/clusters",
			provider.Kubernetes.Cluster.Controller.List)
		kube.GET("/clusters/:id",
			provider.PermissionManager.RequirePermission("kubernetes:cluster:get"),
			provider.Kubernetes.Cluster.Controller.Get)
		kube.POST("/clusters",
			provider.PermissionManager.RequirePermission("kubernetes:cluster:create"),
			provider.Kubernetes.Cluster.Controller.Create)
		kube.PATCH("/clusters/:id",
			provider.PermissionManager.RequirePermission("kubernetes:cluster:update"),
			provider.Kubernetes.Cluster.Controller.Update)
		kube.PATCH("/clusters/:id/token",
			provider.PermissionManager.RequirePermission("kubernetes:cluster:update"),
			provider.Kubernetes.Cluster.Controller.UpdateToken)
		kube.DELETE("/clusters/:id",
			provider.PermissionManager.RequirePermission("kubernetes:cluster:delete"),
			provider.Kubernetes.Cluster.Controller.Delete)
		// Node Routes
		kube.GET("/clusters/:id/nodes",
			provider.Kubernetes.Node.Controller.List)
		kube.GET("/clusters/:id/nodes/:name/detail",
			provider.Kubernetes.Node.Controller.GetDetail)
		kube.PATCH("/clusters/:id/nodes/:name/labels",
			provider.Kubernetes.Node.Controller.Labels)
		kube.PATCH("/clusters/:id/nodes/:name/taints",
			provider.Kubernetes.Node.Controller.Taints)
		kube.PATCH("/clusters/:id/nodes/:name/cordon",
			provider.Kubernetes.Node.Controller.Cordon)
		// Cluster-level Pod list (fieldSelector, labelSelector)
		kube.GET("/clusters/:id/pods",
			provider.Kubernetes.Pod.Controller.ListAll)
		// Namespace Routes
		kube.GET("/clusters/:id/namespaces", provider.Kubernetes.Namespace.Controller.List)
		kube.POST("/clusters/:id/namespaces",
			provider.Kubernetes.Namespace.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/batch-delete",
			provider.Kubernetes.Namespace.Controller.DeleteBatch)
		kube.PATCH("/clusters/:id/namespaces/:namespace/labels",
			provider.Kubernetes.Namespace.Controller.Labels)
		kube.DELETE("/clusters/:id/namespaces/:namespace",
			provider.Kubernetes.Namespace.Controller.Delete)
		// Deployment Routes
		kube.GET("/clusters/:id/namespaces/:namespace/deployments",
			provider.Kubernetes.Deployment.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/deployments",
			provider.Kubernetes.Deployment.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/deployments/batch-delete",
			provider.Kubernetes.Deployment.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/deployments/:name",
			provider.Kubernetes.Deployment.Controller.Get)
		kube.GET("/clusters/:id/namespaces/:namespace/deployments/:name/detail",
			provider.Kubernetes.Deployment.Controller.GetDetail)
		kube.PATCH("/clusters/:id/namespaces/:namespace/deployments/:name",
			provider.Kubernetes.Deployment.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/deployments/:name",
			provider.Kubernetes.Deployment.Controller.Delete)
		kube.PATCH("/clusters/:id/namespaces/:namespace/deployments/:name/scale",
			provider.Kubernetes.Deployment.Controller.Scale)
		kube.POST("/clusters/:id/namespaces/:namespace/deployments/:name/restart",
			provider.Kubernetes.Deployment.Controller.Restart)
		kube.PATCH("/clusters/:id/namespaces/:namespace/deployments/:name/rollout",
			provider.Kubernetes.Deployment.Controller.Rollout)
		// Pod Routes
		kube.GET("/clusters/:id/namespaces/:namespace/pods",
			provider.Kubernetes.Pod.Controller.List)
		kube.DELETE("/clusters/:id/namespaces/:namespace/pods/batch-delete",
			provider.Kubernetes.Pod.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/pods/:name",
			provider.Kubernetes.Pod.Controller.Get)
		kube.GET("/clusters/:id/namespaces/:namespace/pods/:name/detail",
			provider.Kubernetes.Pod.Controller.GetDetail)
		kube.DELETE("/clusters/:id/namespaces/:namespace/pods/:name",
			provider.Kubernetes.Pod.Controller.Delete)
		kube.GET("/clusters/:id/namespaces/:namespace/pods/:name/containers/:container/exec",
			provider.Kubernetes.Pod.Controller.Executor)
		kube.GET("/clusters/:id/namespaces/:namespace/pods/:name/containers/:container/logs",
			provider.Kubernetes.Pod.Controller.GetLogs)
		// HPA Routes
		kube.GET("/clusters/:id/namespaces/:namespace/hpa",
			provider.Kubernetes.Hpa.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/hpa",
			provider.Kubernetes.Hpa.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/hpa/batch-delete",
			provider.Kubernetes.Hpa.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/hpa/:name",
			provider.Kubernetes.Hpa.Controller.Get)
		kube.GET("/clusters/:id/namespaces/:namespace/hpa/:name/history",
			provider.Kubernetes.Hpa.Controller.ListHistory)
		kube.PATCH("/clusters/:id/namespaces/:namespace/hpa/:name",
			provider.Kubernetes.Hpa.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/hpa/:name",
			provider.Kubernetes.Hpa.Controller.Delete)
		// ConfigMap Routes
		kube.GET("/clusters/:id/namespaces/:namespace/configmaps",
			provider.Kubernetes.Configmap.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/configmaps",
			provider.Kubernetes.Configmap.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/configmaps/batch-delete",
			provider.Kubernetes.Configmap.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/configmaps/:name",
			provider.Kubernetes.Configmap.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/configmaps/:name",
			provider.Kubernetes.Configmap.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/configmaps/:name",
			provider.Kubernetes.Configmap.Controller.Delete)
		// StatefulSet Routes
		kube.GET("/clusters/:id/namespaces/:namespace/statefulsets",
			provider.Kubernetes.StatefulSet.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/statefulsets",
			provider.Kubernetes.StatefulSet.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/statefulsets/batch-delete",
			provider.Kubernetes.StatefulSet.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/statefulsets/:name",
			provider.Kubernetes.StatefulSet.Controller.Get)
		kube.GET("/clusters/:id/namespaces/:namespace/statefulsets/:name/detail",
			provider.Kubernetes.StatefulSet.Controller.GetDetail)
		kube.PATCH("/clusters/:id/namespaces/:namespace/statefulsets/:name",
			provider.Kubernetes.StatefulSet.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/statefulsets/:name",
			provider.Kubernetes.StatefulSet.Controller.Delete)
		kube.PATCH("/clusters/:id/namespaces/:namespace/statefulsets/:name/scale",
			provider.Kubernetes.StatefulSet.Controller.Scale)
		kube.POST("/clusters/:id/namespaces/:namespace/statefulsets/:name/restart",
			provider.Kubernetes.StatefulSet.Controller.Restart)
		// DaemonSet Routes
		kube.GET("/clusters/:id/namespaces/:namespace/daemonsets",
			provider.Kubernetes.DaemonSet.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/daemonsets",
			provider.Kubernetes.DaemonSet.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/daemonsets/batch-delete",
			provider.Kubernetes.DaemonSet.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/daemonsets/:name",
			provider.Kubernetes.DaemonSet.Controller.Get)
		kube.GET("/clusters/:id/namespaces/:namespace/daemonsets/:name/detail",
			provider.Kubernetes.DaemonSet.Controller.GetDetail)
		kube.PATCH("/clusters/:id/namespaces/:namespace/daemonsets/:name",
			provider.Kubernetes.DaemonSet.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/daemonsets/:name",
			provider.Kubernetes.DaemonSet.Controller.Delete)
		// ReplicaSet Routes
		kube.GET("/clusters/:id/namespaces/:namespace/replicasets",
			provider.Kubernetes.ReplicaSet.Controller.List)
		kube.DELETE("/clusters/:id/namespaces/:namespace/replicasets/batch-delete",
			provider.Kubernetes.ReplicaSet.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/replicasets/:name",
			provider.Kubernetes.ReplicaSet.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/replicasets/:name",
			provider.Kubernetes.ReplicaSet.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/replicasets/:name",
			provider.Kubernetes.ReplicaSet.Controller.Delete)
		// Service Routes
		kube.GET("/clusters/:id/namespaces/:namespace/services",
			provider.Kubernetes.Service.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/services",
			provider.Kubernetes.Service.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/services/batch-delete",
			provider.Kubernetes.Service.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/services/:name",
			provider.Kubernetes.Service.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/services/:name",
			provider.Kubernetes.Service.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/services/:name",
			provider.Kubernetes.Service.Controller.Delete)
		// Secret Routes
		kube.GET("/clusters/:id/namespaces/:namespace/secrets",
			provider.Kubernetes.Secret.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/secrets",
			provider.Kubernetes.Secret.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/secrets/batch-delete",
			provider.Kubernetes.Secret.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/secrets/:name",
			provider.Kubernetes.Secret.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/secrets/:name",
			provider.Kubernetes.Secret.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/secrets/:name",
			provider.Kubernetes.Secret.Controller.Delete)
		// Role Routes
		kube.GET("/clusters/:id/namespaces/:namespace/roles",
			provider.Kubernetes.Role.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/roles",
			provider.Kubernetes.Role.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/roles/batch-delete",
			provider.Kubernetes.Role.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/roles/:name",
			provider.Kubernetes.Role.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/roles/:name",
			provider.Kubernetes.Role.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/roles/:name",
			provider.Kubernetes.Role.Controller.Delete)
		// RoleBinding Routes
		kube.GET("/clusters/:id/namespaces/:namespace/rolebindings",
			provider.Kubernetes.RoleBinding.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/rolebindings",
			provider.Kubernetes.RoleBinding.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/rolebindings/batch-delete",
			provider.Kubernetes.RoleBinding.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/rolebindings/:name",
			provider.Kubernetes.RoleBinding.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/rolebindings/:name",
			provider.Kubernetes.RoleBinding.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/rolebindings/:name",
			provider.Kubernetes.RoleBinding.Controller.Delete)
		// ClusterRole Routes
		kube.GET("/clusters/:id/clusterroles",
			provider.Kubernetes.ClusterRole.Controller.List)
		kube.POST("/clusters/:id/clusterroles",
			provider.Kubernetes.ClusterRole.Controller.Create)
		kube.DELETE("/clusters/:id/clusterroles/batch-delete",
			provider.Kubernetes.ClusterRole.Controller.DeleteBatch)
		kube.GET("/clusters/:id/clusterroles/:name",
			provider.Kubernetes.ClusterRole.Controller.Get)
		kube.PATCH("/clusters/:id/clusterroles/:name",
			provider.Kubernetes.ClusterRole.Controller.Update)
		kube.DELETE("/clusters/:id/clusterroles/:name",
			provider.Kubernetes.ClusterRole.Controller.Delete)
		// ClusterRoleBinding Routes
		kube.GET("/clusters/:id/clusterrolebindings",
			provider.Kubernetes.ClusterRoleBinding.Controller.List)
		kube.POST("/clusters/:id/clusterrolebindings",
			provider.Kubernetes.ClusterRoleBinding.Controller.Create)
		kube.DELETE("/clusters/:id/clusterrolebindings/batch-delete",
			provider.Kubernetes.ClusterRoleBinding.Controller.DeleteBatch)
		kube.GET("/clusters/:id/clusterrolebindings/:name",
			provider.Kubernetes.ClusterRoleBinding.Controller.Get)
		kube.PATCH("/clusters/:id/clusterrolebindings/:name",
			provider.Kubernetes.ClusterRoleBinding.Controller.Update)
		kube.DELETE("/clusters/:id/clusterrolebindings/:name",
			provider.Kubernetes.ClusterRoleBinding.Controller.Delete)
		// ServiceAccount Routes
		kube.GET("/clusters/:id/namespaces/:namespace/serviceaccounts",
			provider.Kubernetes.ServiceAccount.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/serviceaccounts",
			provider.Kubernetes.ServiceAccount.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/serviceaccounts/batch-delete",
			provider.Kubernetes.ServiceAccount.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/serviceaccounts/:name",
			provider.Kubernetes.ServiceAccount.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/serviceaccounts/:name",
			provider.Kubernetes.ServiceAccount.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/serviceaccounts/:name",
			provider.Kubernetes.ServiceAccount.Controller.Delete)
		// Ingress Routes
		kube.GET("/clusters/:id/namespaces/:namespace/ingresses",
			provider.Kubernetes.Ingress.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/ingresses",
			provider.Kubernetes.Ingress.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/ingresses/batch-delete",
			provider.Kubernetes.Ingress.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/ingresses/:name",
			provider.Kubernetes.Ingress.Controller.Get)
		kube.GET("/clusters/:id/namespaces/:namespace/ingresses/:name/detail",
			provider.Kubernetes.Ingress.Controller.GetDetail)
		kube.PATCH("/clusters/:id/namespaces/:namespace/ingresses/:name",
			provider.Kubernetes.Ingress.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/ingresses/:name",
			provider.Kubernetes.Ingress.Controller.Delete)
		// IngressClass Routes
		kube.GET("/clusters/:id/ingressclasses",
			provider.Kubernetes.IngressClass.Controller.List)
		kube.POST("/clusters/:id/ingressclasses",
			provider.Kubernetes.IngressClass.Controller.Create)
		kube.DELETE("/clusters/:id/ingressclasses/batch-delete",
			provider.Kubernetes.IngressClass.Controller.DeleteBatch)
		kube.GET("/clusters/:id/ingressclasses/:name",
			provider.Kubernetes.IngressClass.Controller.Get)
		kube.PATCH("/clusters/:id/ingressclasses/:name",
			provider.Kubernetes.IngressClass.Controller.Update)
		kube.DELETE("/clusters/:id/ingressclasses/:name",
			provider.Kubernetes.IngressClass.Controller.Delete)
		// StorageClass Routes
		kube.GET("/clusters/:id/storageclass",
			provider.Kubernetes.StorageClass.Controller.List)
		kube.POST("/clusters/:id/storageclass",
			provider.Kubernetes.StorageClass.Controller.Create)
		kube.DELETE("/clusters/:id/storageclass/batch-delete",
			provider.Kubernetes.StorageClass.Controller.DeleteBatch)
		kube.GET("/clusters/:id//storageclass/:name",
			provider.Kubernetes.StorageClass.Controller.Get)
		kube.PATCH("/clusters/:id/storageclass/:name",
			provider.Kubernetes.StorageClass.Controller.Update)
		kube.DELETE("/clusters/:id/storageclass/:name",
			provider.Kubernetes.StorageClass.Controller.Delete)
		// PersistentVolume Routes
		kube.GET("/clusters/:id/persistentvolumes",
			provider.Kubernetes.PersistentVolume.Controller.List)
		kube.POST("/clusters/:id/persistentvolumes",
			provider.Kubernetes.PersistentVolume.Controller.Create)
		kube.DELETE("/clusters/:id/persistentvolumes/batch-delete",
			provider.Kubernetes.PersistentVolume.Controller.DeleteBatch)
		kube.GET("/clusters/:id/persistentvolumes/:name",
			provider.Kubernetes.PersistentVolume.Controller.Get)
		kube.PATCH("/clusters/:id/persistentvolumes/:name",
			provider.Kubernetes.PersistentVolume.Controller.Update)
		kube.DELETE("/clusters/:id/persistentvolumes/:name",
			provider.Kubernetes.PersistentVolume.Controller.Delete)
		// PersistentVolumeClaim Routes
		kube.GET("/clusters/:id/namespaces/:namespace/persistentvolumeclaims",
			provider.Kubernetes.PersistentVolumeClaim.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/persistentvolumeclaims",
			provider.Kubernetes.PersistentVolumeClaim.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/persistentvolumeclaims/batch-delete",
			provider.Kubernetes.PersistentVolumeClaim.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/persistentvolumeclaims/:name",
			provider.Kubernetes.PersistentVolumeClaim.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/persistentvolumeclaims/:name",
			provider.Kubernetes.PersistentVolumeClaim.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/persistentvolumeclaims/:name",
			provider.Kubernetes.PersistentVolumeClaim.Controller.Delete)
		// Event Routes
		kube.GET("/clusters/:id/namespaces/:namespace/events",
			provider.Kubernetes.Event.Controller.List)
		// Task Routes
		kube.GET("/clusters/:id/namespaces/:namespace/tasks",
			provider.Kubernetes.Task.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/tasks",
			provider.Kubernetes.Task.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/tasks/batch-delete",
			provider.Kubernetes.Task.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/tasks/:name",
			provider.Kubernetes.Task.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/tasks/:name",
			provider.Kubernetes.Task.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/tasks/:name",
			provider.Kubernetes.Task.Controller.Delete)
		// TaskRun Routes
		kube.GET("/clusters/:id/namespaces/:namespace/taskruns",
			provider.Kubernetes.TaskRun.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/taskruns",
			provider.Kubernetes.TaskRun.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/taskruns/batch-delete",
			provider.Kubernetes.TaskRun.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/taskruns/:name",
			provider.Kubernetes.TaskRun.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/taskruns/:name",
			provider.Kubernetes.TaskRun.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/taskruns/:name",
			provider.Kubernetes.TaskRun.Controller.Delete)
		// Pipeline Routes
		kube.GET("/clusters/:id/namespaces/:namespace/pipelines",
			provider.Kubernetes.Pipeline.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/pipelines",
			provider.Kubernetes.Pipeline.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/pipelines/batch-delete",
			provider.Kubernetes.Pipeline.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/pipelines/:name",
			provider.Kubernetes.Pipeline.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/pipelines/:name",
			provider.Kubernetes.Pipeline.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/pipelines/:name",
			provider.Kubernetes.Pipeline.Controller.Delete)
		// PipelineRun Routes
		kube.GET("/clusters/:id/namespaces/:namespace/pipelineruns",
			provider.Kubernetes.PipelineRun.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/pipelineruns",
			provider.Kubernetes.PipelineRun.Controller.Create)
		kube.DELETE("/clusters/:id/namespaces/:namespace/pipelineruns/batch-delete",
			provider.Kubernetes.PipelineRun.Controller.DeleteBatch)
		kube.GET("/clusters/:id/namespaces/:namespace/pipelineruns/:name",
			provider.Kubernetes.PipelineRun.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/pipelineruns/:name",
			provider.Kubernetes.PipelineRun.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/pipelineruns/:name",
			provider.Kubernetes.PipelineRun.Controller.Delete)
	}
	// 别名路由：兼容 /api/v1/clusters/:id/token
	rg.PATCH("/clusters/:id/token", provider.Kubernetes.Cluster.Controller.UpdateToken)
}
