package dto

import (
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type ConfigMap struct {
	Name      string `json:"name"`
	Data      int    `json:"data"`
	CreatedAt string `json:"createdAt"`
}

// HPA 水平扩缩列表项
type HPA struct {
	Namespace       string `json:"namespace"`
	Name            string `json:"name"`
	ScaleTargetKind string `json:"scaleTargetKind"`
	ScaleTargetName string `json:"scaleTargetName"`
	MinReplicas     *int32 `json:"minReplicas,omitempty"`
	MaxReplicas     int32  `json:"maxReplicas"`
	CurrentReplicas int32  `json:"currentReplicas"`
	DesiredReplicas int32  `json:"desiredReplicas"`
	CreatedAt       string `json:"createdAt"`
}

type Event struct {
	Type     string `json:"type"`
	Reason   string `json:"reason"`
	Object   string `json:"object"`
	Source   string `json:"source"`
	Message  string `json:"message"`
	CreateAt string `json:"createAt"`
}

type Pod struct {
	Namespace    string                `json:"namespace"`
	Name         string                `json:"name"`
	PodIP        string                `json:"podIP"`
	HostIP       string                `json:"hostIP"`
	NodeName     string                `json:"nodeName"`
	Status       string                `json:"status"`
	RestartCount int                   `json:"restartCount"`
	Ready        string                `json:"ready"`
	Containers   []Container           `json:"containers"`
	Conditions   []corev1.PodCondition `json:"conditions"`
	Labels       map[string]string     `json:"labels"`
	CreatedAt    string                `json:"createdAt"`
}

type Container struct {
	Name         string `json:"name"`
	Image        string `json:"image"`
	ImageID      string `json:"imageID"`
	Ready        bool   `json:"ready"`
	State        string `json:"state"`
	RestartCount int32  `json:"restartCount"`
	StartedAt    string `json:"startedAt,omitempty"`
}

type Node struct {
	Name                    string                 `json:"name,omitempty"`
	Address                 string                 `json:"address,omitempty"`
	OSImage                 string                 `json:"OSImage,omitempty"`
	OSSystem                string                 `json:"OSSystem,omitempty"`
	Architecture            string                 `json:"architecture,omitempty"`
	KubeletVersion          string                 `json:"kubeletVersion,omitempty"`
	KubeProxyVersion        string                 `json:"kubeProxyVersion,omitempty"`
	KernelVersion           string                 `json:"kernelVersion,omitempty"`
	ContainerRuntimeVersion string                 `json:"containerRuntimeVersion,omitempty"`
	PodCIDR                 string                 `json:"podCIDR,omitempty"`
	Labels                  map[string]string      `json:"labels,omitempty"`
	Taints                  []corev1.Taint         `json:"taints,omitempty"`
	Status                  corev1.ConditionStatus `json:"status,omitempty"`
	MachineID               string                 `json:"machineID,omitempty"`
	SystemUUID              string                 `json:"systemUUID,omitempty"`
	BootID                  string                 `json:"bootID,omitempty"`
	Unschedulable           *bool                  `json:"unschedulable,omitempty"`
	Conditions              []NodeCondition        `json:"conditions,omitempty"`
	Pods                    []Pod                  `json:"pods,omitempty"`
	CreatedAt               string                 `json:"createdAt,omitempty"`
}

type NodeCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	Reason             string `json:"reason,omitempty"`
	Message            string `json:"message,omitempty"`
	LastHeartbeatTime  string `json:"lastHeartbeatTime"`
	LastTransitionTime string `json:"lastTransitionTime"`
}

type PersistentVolume struct {
	Name          string                               `json:"name"`
	Capacity      string                               `json:"capacity"`
	AccessModes   []corev1.PersistentVolumeAccessMode  `json:"accessModes"`
	ReclaimPolicy corev1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy"`
	Status        string                               `json:"status"`
	Claim         string                               `json:"claim"`
	StorageClass  string                               `json:"storageClass"`
	Reason        string                               `json:"reason"`
	CreatedAt     string                               `json:"createdAt"`
}

type PersistentVolumeClaim struct {
	Name         string                              `json:"name"`
	Capacity     string                              `json:"capacity"`
	AccessModes  []corev1.PersistentVolumeAccessMode `json:"accessModes"`
	Status       string                              `json:"status"`
	StorageClass *string                             `json:"storageClass"`
	Volume       string                              `json:"volume"`
	CreatedAt    string                              `json:"createdAt"`
}

type Namespace struct {
	Name      string                `json:"name"`
	Status    corev1.NamespacePhase `json:"status"`
	Labels    map[string]string     `json:"labels"`
	CreatedAt string                `json:"createdAt"`
}

type ReplicaSet struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	Revision  string `json:"revision"`
	CreatedAt string `json:"createdAt"`
}

type Secret struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Data      int    `json:"data"`
	CreatedAt string `json:"createdAt"`
}

type Role struct {
	Name      string `json:"name"`
	Rules     int    `json:"rules"`
	CreatedAt string `json:"createdAt"`
}

type ClusterRole struct {
	Name      string `json:"name"`
	Rules     int    `json:"rules"`
	CreatedAt string `json:"createdAt"`
}

type RoleBinding struct {
	Name      string `json:"name"`
	RoleRef   string `json:"roleRef"`
	Subjects  int    `json:"subjects"`
	CreatedAt string `json:"createdAt"`
}

type ClusterRoleBinding struct {
	Name      string `json:"name"`
	RoleRef   string `json:"roleRef"`
	Subjects  int    `json:"subjects"`
	CreatedAt string `json:"createdAt"`
}

type ServiceAccount struct {
	Name             string `json:"name"`
	Secrets          int    `json:"secrets"`
	ImagePullSecrets int    `json:"imagePullSecrets"`
	CreatedAt        string `json:"createdAt"`
}

type Ingress struct {
	Name      string   `json:"name"`
	Hosts     []string `json:"hosts"`
	Paths     int      `json:"paths"`
	TLS       bool     `json:"tls"`
	Endpoint  string   `json:"endpoint"`
	ClassName string   `json:"className"`
	CreatedAt string   `json:"createdAt"`
}

type IngressDetail struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	ClassName   string            `json:"className"`
	Rules       []IngressRule     `json:"rules"`
	Annotations map[string]string `json:"annotations"`
	Labels      map[string]string `json:"labels"`
	CreatedAt   string            `json:"createdAt"`
}

type IngressRule struct {
	Host  string            `json:"host"`
	Paths []IngressRulePath `json:"paths"`
}

type IngressRulePath struct {
	Path     string `json:"path"`
	PathType string `json:"pathType"`
	Service  string `json:"service"`
	Port     int    `json:"port"`
}

type IngressClass struct {
	Name       string `json:"name"`
	Controller string `json:"controller"`
	Parameters string `json:"parameters"`
	CreatedAt  string `json:"createdAt"`
}

type StorageClass struct {
	Name                 string `json:"name"`
	Provisioner          string `json:"provisioner"`
	ReclaimPolicy        string `json:"reclaimPolicy"`
	VolumeBindingMode    string `json:"volumeBindingMode"`
	AllowVolumeExpansion bool   `json:"allowVolumeExpansion"`
	CreatedAt            string `json:"createdAt"`
}

type StatefulSet struct {
	Namespace         string                           `json:"namespace"`
	Name              string                           `json:"name"`
	Replicas          int32                            `json:"replicas"`
	ReadyReplicas     int32                            `json:"readyReplicas"`
	CurrentReplicas   int32                            `json:"currentReplicas"`
	AvailableReplicas int32                            `json:"availableReplicas"`
	UpdatedReplicas   int32                            `json:"updatedReplicas"`
	MatchLabels       map[string]string                `json:"matchLabels"`
	Strategy          appsv1.StatefulSetUpdateStrategy `json:"strategy"`
	CreatedAt         string                           `json:"createdAt"`
}

type DaemonSet struct {
	Namespace              string                         `json:"namespace"`
	Name                   string                         `json:"name"`
	CurrentNumberScheduled int32                          `json:"currentNumberScheduled"`
	DesiredNumberScheduled int32                          `json:"desiredNumberScheduled"`
	NumberReady            int32                          `json:"numberReady"`
	NumberAvailable        int32                          `json:"numberAvailable"`
	UpdatedNumberScheduled int32                          `json:"updatedNumberScheduled"`
	MatchLabels            map[string]string              `json:"matchLabels"`
	Strategy               appsv1.DaemonSetUpdateStrategy `json:"strategy"`
	CreatedAt              string                         `json:"createdAt"`
}

type Deployment struct {
	Namespace         string                       `json:"namespace"`
	Name              string                       `json:"name"`
	Replicas          int32                        `json:"replicas"`
	ReadyReplicas     int32                        `json:"readyReplicas"`
	AvailableReplicas int32                        `json:"availableReplicas"`
	UpdatedReplicas   int32                        `json:"updatedReplicas"`
	Paused            bool                         `json:"paused"`
	MatchLabels       map[string]string            `json:"matchLabels"`
	Strategy          appsv1.DeploymentStrategy    `json:"strategy"`
	Conditions        []appsv1.DeploymentCondition `json:"conditions"`
	CreatedAt         string                       `json:"createdAt"`
	UpdatedAt         string                       `json:"updatedAt,omitempty"`
}

type Service struct {
	Name      string             `json:"name,omitempty"`
	Type      corev1.ServiceType `json:"type,omitempty"`
	ClusterIP string             `json:"clusterIP,omitempty"`
	Ports     string             `json:"ports,omitempty"`
	CreatedAt string             `json:"createdAt,omitempty"`
}

type Job struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type CronJob struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type Pipeline struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

type PipelineRun struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

type Task struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

type TaskRun struct {
	Name      string                 `json:"name,omitempty"`
	Namespace string                 `json:"namespace,omitempty"`
	Status    corev1.ConditionStatus `json:"status"`
	Reason    string                 `json:"reason,omitempty"`
	StartAt   string                 `json:"startAt,omitempty"`
	EndAt     string                 `json:"endAt,omitempty"`
	CreatedAt string                 `json:"createdAt,omitempty"`
}

type ClusterDTO struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Host        string    `json:"host,omitempty"`
	Version     string    `json:"version,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

// PermissionDTO RBAC 权限 DTO
type PermissionDTO struct {
	ID          uint       `json:"id,omitempty"`
	UserID      uint       `json:"userID,omitempty"`
	Username    string     `json:"username,omitempty"`
	ClusterID   uint       `json:"clusterID,omitempty"`
	ClusterName string     `json:"clusterName,omitempty"`
	Namespace   string     `json:"namespace,omitempty"`
	Role        string     `json:"role,omitempty"`
	Creator     string     `json:"creator,omitempty"`
	CreatedAt   time.Time  `json:"createdAt,omitempty"`
	UpdatedAt   time.Time  `json:"updatedAt,omitempty"`
	SyncStatus  string     `json:"syncStatus,omitempty"`
	RetryCount  int        `json:"retryCount,omitempty"`
	LastError   string     `json:"lastError,omitempty"`
	LastErrorAt *time.Time `json:"lastErrorAt,omitempty"`
	LastSyncAt  *time.Time `json:"lastSyncAt,omitempty"`
	NextRetryAt *time.Time `json:"nextRetryAt,omitempty"`
}
