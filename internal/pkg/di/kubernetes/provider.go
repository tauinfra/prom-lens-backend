package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/pkg/encryption"
	"valyria-backend/internal/pkg/k8s/factory"

	"gorm.io/gorm"
)

type KubernetesProvider struct {
	Policy                *PolicyProvider
	Cluster               *ClusterProvider
	Node                  *NodeProvider
	Namespace             *NamespaceProvider
	Deployment            *DeploymentProvider
	Pod                   *PodProvider
	Event                 *EventProvider
	Configmap             *ConfigmapProvider
	Service               *ServiceProvider
	Secret                *SecretProvider
	DaemonSet             *DaemonSetProvider
	ReplicaSet            *ReplicaSetProvider
	StatefulSet           *StatefulSetProvider
	StorageClass          *StorageClassProvider
	PersistentVolume      *PersistentVolumeProvider
	PersistentVolumeClaim *PersistentVolumeClaimProvider
	Task                  *TaskProvider
	TaskRun               *TaskRunProvider
	Pipeline              *PipelineProvider
	PipelineRun           *PipelineRunProvider
	// ... 其他资源
}

func NewKubernetesProvider(db *gorm.DB, encryptor encryption.Encryptor) *KubernetesProvider {
	// 先创建基础依赖
	repo := repository.NewClusterRepository(db, encryptor)

	// 创建共享的工厂
	cfgFactory := repository.NewKubeConfigFactory(repo, encryptor)
	gvkFactory := factory.NewGVKFactory()
	return &KubernetesProvider{
		Policy:                NewPolicyProvider(db),
		Cluster:               NewClusterProvider(db, encryptor),
		Node:                  NewNodeProvider(cfgFactory, gvkFactory),
		Namespace:             NewNamespaceProvider(cfgFactory, gvkFactory),
		Deployment:            NewDeploymentProvider(cfgFactory, gvkFactory),
		Pod:                   NewPodProvider(cfgFactory, gvkFactory),
		Event:                 NewEventProvider(cfgFactory),
		Configmap:             NewConfigmapProvider(cfgFactory, gvkFactory),
		Service:               NewServiceProvider(cfgFactory, gvkFactory),
		Secret:                NewSecretProvider(cfgFactory, gvkFactory),
		DaemonSet:             NewRDaemonSetProvider(cfgFactory, gvkFactory),
		ReplicaSet:            NewReplicaSetProvider(cfgFactory, gvkFactory),
		StatefulSet:           NewStatefulSetProvider(cfgFactory, gvkFactory),
		StorageClass:          NewStorageClassProvider(cfgFactory, gvkFactory),
		PersistentVolume:      NewPersistentVolumeProvider(cfgFactory, gvkFactory),
		PersistentVolumeClaim: NewPersistentVolumeClaimProvider(cfgFactory, gvkFactory),
		Task:                  NewTaskProvider(cfgFactory, gvkFactory),
		TaskRun:               NewTaskRunProvider(cfgFactory, gvkFactory),
		Pipeline:              NewPipelineProvider(cfgFactory, gvkFactory),
		PipelineRun:           NewPipelineRunProvider(cfgFactory, gvkFactory),
	}
}
