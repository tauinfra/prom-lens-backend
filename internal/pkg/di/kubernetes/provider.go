package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/core/config"
	"valyria-backend/internal/pkg/encryption"
	"valyria-backend/internal/pkg/k8s/factory"

	"gorm.io/gorm"
)

type KubernetesProvider struct {
	CfgFactory             *repository.KubeConfigFactory
	GvkFactory             *factory.GVKFactory
	Cluster                *ClusterProvider
	Permission             *PermissionProvider
	Node                   *NodeProvider
	Namespace              *NamespaceProvider
	Deployment             *DeploymentProvider
	Pod                    *PodProvider
	Event                 *EventProvider
	Configmap             *ConfigmapProvider
	Service               *ServiceProvider
	Secret                *SecretProvider
	Role                  *RoleProvider
	ClusterRole           *ClusterRoleProvider
	RoleBinding           *RoleBindingProvider
	ClusterRoleBinding    *ClusterRoleBindingProvider
	ServiceAccount        *ServiceAccountProvider
	Ingress               *IngressProvider
	IngressClass          *IngressClassProvider
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
	Hpa                   *HpaProvider
	HpaHistorySvc         service.HpaHistoryManager // 供 HpaHistorySyncWorker 使用
	// ... 其他资源
}

func NewKubernetesProvider(db *gorm.DB, encryptor encryption.Encryptor, cfg *config.Config) *KubernetesProvider {
	// 先创建基础依赖
	repo := repository.NewClusterRepository(db, encryptor)

	// 创建共享的工厂
	cfgFactory := repository.NewKubeConfigFactory(repo, encryptor)
	gvkFactory := factory.NewGVKFactory()

	// HPA 扩缩容历史：落库 + 定时从 Event 同步
	hpaHistoryRepo := repository.NewHpaHistoryRepository(db)
	hpaHistorySvc := service.NewHpaHistoryManager(hpaHistoryRepo, cfgFactory, repo)

	return &KubernetesProvider{
		CfgFactory:             cfgFactory,
		GvkFactory:             gvkFactory,
		Cluster:                NewClusterProvider(db, encryptor),
		Permission:             NewPermissionProvider(db, cfgFactory, gvkFactory, cfg),
		Node:                   NewNodeProvider(cfgFactory, gvkFactory),
		Namespace:              NewNamespaceProvider(cfgFactory, gvkFactory),
		Deployment:             NewDeploymentProvider(cfgFactory, gvkFactory),
		Pod:                    NewPodProvider(cfgFactory, gvkFactory),
		Event:                 NewEventProvider(cfgFactory),
		Configmap:             NewConfigmapProvider(cfgFactory, gvkFactory),
		Service:               NewServiceProvider(cfgFactory, gvkFactory),
		Secret:                NewSecretProvider(cfgFactory, gvkFactory),
		Role:                  NewRoleProvider(cfgFactory, gvkFactory),
		ClusterRole:           NewClusterRoleProvider(cfgFactory, gvkFactory),
		RoleBinding:           NewRoleBindingProvider(cfgFactory, gvkFactory),
		ClusterRoleBinding:    NewClusterRoleBindingProvider(cfgFactory, gvkFactory),
		ServiceAccount:        NewServiceAccountProvider(cfgFactory, gvkFactory),
		Ingress:               NewIngressProvider(cfgFactory, gvkFactory),
		IngressClass:          NewIngressClassProvider(cfgFactory, gvkFactory),
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
		Hpa:                   NewHpaProvider(cfgFactory, gvkFactory, hpaHistorySvc),
		HpaHistorySvc:         hpaHistorySvc,
	}
}
