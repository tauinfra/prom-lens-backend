package routes

import (
	"valyria-backend/internal/pkg/di"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func kubeRouters(rg *gin.RouterGroup, tx *gorm.DB, provider *di.Provider) {
	kube := rg.Group("/kubernetes")
	kube.Use(provider.PermissionManager.RBACPermission(tx))
	{
		// Cluster Routes
		kube.GET("/polices", provider.Kubernetes.Policy.Controller.List)
		kube.POST("/polices", provider.Kubernetes.Policy.Controller.Create)
		kube.DELETE("/polices/:id", provider.Kubernetes.Policy.Controller.Delete)
		// Cluster Routes
		kube.GET("/clusters", provider.Kubernetes.Cluster.Controller.List)
		kube.GET("/clusters/:id", provider.Kubernetes.Cluster.Controller.Get)
		kube.POST("/clusters", provider.Kubernetes.Cluster.Controller.Create)
		kube.PATCH("/clusters/:id", provider.Kubernetes.Cluster.Controller.Update)
		kube.DELETE("/clusters/:id", provider.Kubernetes.Cluster.Controller.Delete)
		// Node Routes
		kube.GET("/clusters/:id/nodes", provider.Kubernetes.Node.Controller.List)
		kube.GET("/clusters/:id/nodes/:name/detail", provider.Kubernetes.Node.Controller.GetDetail)
		kube.PATCH("/clusters/:id/nodes/:name/labels", provider.Kubernetes.Node.Controller.Labels)
		kube.PATCH("/clusters/:id/nodes/:name/taints", provider.Kubernetes.Node.Controller.Taints)
		kube.PATCH("/clusters/:id/nodes/:name/cordon", provider.Kubernetes.Node.Controller.Cordon)
		// Namespace Routes
		kube.GET("/clusters/:id/namespaces", provider.Kubernetes.Namespace.Controller.List)
		kube.POST("/clusters/:id/namespaces", provider.Kubernetes.Namespace.Controller.Create)
		// Deployment Routes
		kube.GET("/clusters/:id/namespaces/:namespace/deployments", provider.Kubernetes.Deployment.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/deployments", provider.Kubernetes.Deployment.Controller.Create)
		kube.GET("/clusters/:id/namespaces/:namespace/deployments/:name", provider.Kubernetes.Deployment.Controller.Get)
		kube.GET("/clusters/:id/namespaces/:namespace/deployments/:name/detail", provider.Kubernetes.Deployment.Controller.GetDetail)
		kube.PATCH("/clusters/:id/namespaces/:namespace/deployments/:name", provider.Kubernetes.Deployment.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/deployments/:name", provider.Kubernetes.Deployment.Controller.Delete)
		kube.PATCH("/clusters/:id/namespaces/:namespace/deployments/:name/scale", provider.Kubernetes.Deployment.Controller.Scale)
		kube.POST("/clusters/:id/namespaces/:namespace/deployments/:name/restart", provider.Kubernetes.Deployment.Controller.Restart)
		kube.PATCH("/clusters/:id/namespaces/:namespace/deployments/:name/rollout", provider.Kubernetes.Deployment.Controller.Rollout)
		// Pod Routes
		kube.GET("/clusters/:id/namespaces/:namespace/pods", provider.Kubernetes.Pod.Controller.List)
		kube.GET("/clusters/:id/namespaces/:namespace/pods/:name", provider.Kubernetes.Pod.Controller.Get)
		kube.GET("/clusters/:id/namespaces/:namespace/pods/:name/detail", provider.Kubernetes.Pod.Controller.GetDetail)
		kube.GET("/clusters/:id/namespaces/:namespace/pods/:name/containers/:container/exec", provider.Kubernetes.Pod.Controller.Executor)
		kube.GET("/clusters/:id/namespaces/:namespace/pods/:name/containers/:container/logs", provider.Kubernetes.Pod.Controller.GetLogs)
		// ConfigMap Routes
		kube.GET("/clusters/:id/namespaces/:namespace/configmaps", provider.Kubernetes.Configmap.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/configmaps", provider.Kubernetes.Configmap.Controller.Create)
		kube.GET("/clusters/:id/namespaces/:namespace/configmaps/:name", provider.Kubernetes.Configmap.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/configmaps/:name", provider.Kubernetes.Configmap.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/configmaps/:name", provider.Kubernetes.Configmap.Controller.Delete)
		// StatefulSet Routes
		kube.GET("/clusters/:id/namespaces/:namespace/statefulsets", provider.Kubernetes.StatefulSet.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/statefulsets", provider.Kubernetes.StatefulSet.Controller.Create)
		kube.GET("/clusters/:id/namespaces/:namespace/statefulsets/:name", provider.Kubernetes.StatefulSet.Controller.Get)
		kube.GET("/clusters/:id/namespaces/:namespace/statefulsets/:name/detail", provider.Kubernetes.StatefulSet.Controller.GetDetail)
		kube.PATCH("/clusters/:id/namespaces/:namespace/statefulsets/:name", provider.Kubernetes.StatefulSet.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/statefulsets/:name", provider.Kubernetes.StatefulSet.Controller.Delete)
		kube.PATCH("/clusters/:id/namespaces/:namespace/statefulsets/:name/scale", provider.Kubernetes.StatefulSet.Controller.Scale)
		kube.POST("/clusters/:id/namespaces/:namespace/statefulsets/:name/restart", provider.Kubernetes.StatefulSet.Controller.Restart)
		// DaemonSet Routes
		kube.GET("/clusters/:id/namespaces/:namespace/daemonsets", provider.Kubernetes.DaemonSet.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/daemonsets", provider.Kubernetes.DaemonSet.Controller.Create)
		kube.GET("/clusters/:id/namespaces/:namespace/daemonsets/:name", provider.Kubernetes.DaemonSet.Controller.Get)
		kube.GET("/clusters/:id/namespaces/:namespace/daemonsets/:name/detail", provider.Kubernetes.DaemonSet.Controller.GetDetail)
		kube.PATCH("/clusters/:id/namespaces/:namespace/daemonsets/:name", provider.Kubernetes.DaemonSet.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/daemonsets/:name", provider.Kubernetes.DaemonSet.Controller.Delete)
		// ReplicaSet Routes
		kube.GET("/clusters/:id/namespaces/:namespace/replicasets", provider.Kubernetes.ReplicaSet.Controller.List)
		kube.GET("/clusters/:id/namespaces/:namespace/replicasets/:name", provider.Kubernetes.ReplicaSet.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/replicasets/:name", provider.Kubernetes.ReplicaSet.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/replicasets/:name", provider.Kubernetes.ReplicaSet.Controller.Delete)
		// Service Routes
		kube.GET("/clusters/:id/namespaces/:namespace/services", provider.Kubernetes.Service.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/services", provider.Kubernetes.Service.Controller.Create)
		kube.GET("/clusters/:id/namespaces/:namespace/services/:name", provider.Kubernetes.Service.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/services/:name", provider.Kubernetes.Service.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/services/:name", provider.Kubernetes.Service.Controller.Delete)
		// Secret Routes
		kube.GET("/clusters/:id/namespaces/:namespace/secrets", provider.Kubernetes.Secret.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/secrets", provider.Kubernetes.Secret.Controller.Create)
		kube.GET("/clusters/:id/namespaces/:namespace/secrets/:name", provider.Kubernetes.Secret.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/secrets/:name", provider.Kubernetes.Secret.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/secrets/:name", provider.Kubernetes.Secret.Controller.Delete)
		// StorageClass Routes
		kube.GET("/clusters/:id/storageclass", provider.Kubernetes.StorageClass.Controller.List)
		kube.POST("/clusters/:id/storageclass", provider.Kubernetes.StorageClass.Controller.Create)
		kube.GET("/clusters/:id//storageclass/:name", provider.Kubernetes.StorageClass.Controller.Get)
		kube.PATCH("/clusters/:id/storageclass/:name", provider.Kubernetes.StorageClass.Controller.Update)
		kube.DELETE("/clusters/:id/storageclass/:name", provider.Kubernetes.StorageClass.Controller.Delete)
		// PersistentVolume Routes
		kube.GET("/clusters/:id/persistentvolumes", provider.Kubernetes.PersistentVolume.Controller.List)
		kube.POST("/clusters/:id/persistentvolumes", provider.Kubernetes.PersistentVolume.Controller.Create)
		kube.GET("/clusters/:id/persistentvolumes/:name", provider.Kubernetes.PersistentVolume.Controller.Get)
		kube.PATCH("/clusters/:id/persistentvolumes/:name", provider.Kubernetes.PersistentVolume.Controller.Update)
		kube.DELETE("/clusters/:id/persistentvolumes/:name", provider.Kubernetes.PersistentVolume.Controller.Delete)
		// PersistentVolumeClaim Routes
		kube.GET("/clusters/:id/namespaces/:namespace/persistentvolumeclaims", provider.Kubernetes.PersistentVolumeClaim.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/persistentvolumeclaims", provider.Kubernetes.PersistentVolumeClaim.Controller.Create)
		kube.GET("/clusters/:id/namespaces/:namespace/persistentvolumeclaims/:name", provider.Kubernetes.PersistentVolumeClaim.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/persistentvolumeclaims/:name", provider.Kubernetes.PersistentVolumeClaim.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/persistentvolumeclaims/:name", provider.Kubernetes.PersistentVolumeClaim.Controller.Delete)
		// Event Routes
		kube.GET("/clusters/:id/namespaces/:namespace/events", provider.Kubernetes.Event.Controller.List)
		// Task Routes
		kube.GET("/clusters/:id/namespaces/:namespace/tasks", provider.Kubernetes.Task.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/tasks", provider.Kubernetes.Task.Controller.Create)
		kube.GET("/clusters/:id/namespaces/:namespace/tasks/:name", provider.Kubernetes.Task.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/tasks/:name", provider.Kubernetes.Task.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/tasks/:name", provider.Kubernetes.Task.Controller.Delete)
		// TaskRun Routes
		kube.GET("/clusters/:id/namespaces/:namespace/taskruns", provider.Kubernetes.TaskRun.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/taskruns", provider.Kubernetes.TaskRun.Controller.Create)
		kube.GET("/clusters/:id/namespaces/:namespace/taskruns/:name", provider.Kubernetes.TaskRun.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/taskruns/:name", provider.Kubernetes.TaskRun.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/taskruns/:name", provider.Kubernetes.TaskRun.Controller.Delete)
		// Pipeline Routes
		kube.GET("/clusters/:id/namespaces/:namespace/pipelines", provider.Kubernetes.Pipeline.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/pipelines", provider.Kubernetes.Pipeline.Controller.Create)
		kube.GET("/clusters/:id/namespaces/:namespace/pipelines/:name", provider.Kubernetes.Pipeline.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/pipelines/:name", provider.Kubernetes.Pipeline.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/pipelines/:name", provider.Kubernetes.Pipeline.Controller.Delete)
		// PipelineRun Routes
		kube.GET("/clusters/:id/namespaces/:namespace/pipelineruns", provider.Kubernetes.PipelineRun.Controller.List)
		kube.POST("/clusters/:id/namespaces/:namespace/pipelineruns", provider.Kubernetes.PipelineRun.Controller.Create)
		kube.GET("/clusters/:id/namespaces/:namespace/pipelineruns/:name", provider.Kubernetes.PipelineRun.Controller.Get)
		kube.PATCH("/clusters/:id/namespaces/:namespace/pipelineruns/:name", provider.Kubernetes.PipelineRun.Controller.Update)
		kube.DELETE("/clusters/:id/namespaces/:namespace/pipelineruns/:name", provider.Kubernetes.PipelineRun.Controller.Delete)
	}
}
