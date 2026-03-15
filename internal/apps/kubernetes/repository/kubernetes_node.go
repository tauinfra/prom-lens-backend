package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/factory"
	"valyria-backend/internal/pkg/k8s/helper"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeRepository interface {
	List(ctx context.Context, id uint) ([]Node, error)
	Get(ctx context.Context, id uint, name string) (*corev1.Node, error)
	GetDetail(ctx context.Context, id uint, name string) (*Node, error)
	Update(ctx context.Context, id uint, node *corev1.Node) (*corev1.Node, error)
	Cordon(ctx context.Context, id uint, name string) (*corev1.Node, error) // 调度
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
	Conditions              []helper.NodeCondition `json:"conditions,omitempty"`
	Pods                    []Pod                  `json:"pods,omitempty"`
	CreatedAt               string                 `json:"createdAt,omitempty"`
}

type nodes struct {
	cfgFactory *KubeConfigFactory
	gvkFactory *factory.GVKFactory
}

func NewNodeRepository(cfgFactory *KubeConfigFactory, gvkFactory *factory.GVKFactory) NodeRepository {
	return &nodes{cfgFactory: cfgFactory, gvkFactory: gvkFactory}
}

func (r *nodes) List(ctx context.Context, id uint) ([]Node, error) {
	var (
		data []Node
		node Node
	)

	nodeHelper := helper.NewNodeHelper()

	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return data, fmt.Errorf("kubernetes CoreV1 nodes list failed. err: %v", err)
	}

	for _, item := range response.Items {
		node = Node{
			Name:                    item.Name,
			Labels:                  item.Labels,
			KubeletVersion:          item.Status.NodeInfo.KubeletVersion,
			KubeProxyVersion:        item.Status.NodeInfo.KubeProxyVersion,
			KernelVersion:           item.Status.NodeInfo.KernelVersion,
			ContainerRuntimeVersion: item.Status.NodeInfo.ContainerRuntimeVersion,
			OSSystem:                item.Status.NodeInfo.OperatingSystem,
			OSImage:                 item.Status.NodeInfo.OSImage,
			Architecture:            item.Status.NodeInfo.Architecture,
			MachineID:               item.Status.NodeInfo.MachineID,
			SystemUUID:              item.Status.NodeInfo.SystemUUID,
			BootID:                  item.Status.NodeInfo.BootID,
			PodCIDR:                 item.Spec.PodCIDR,
			Taints:                  item.Spec.Taints,
			Unschedulable:           &item.Spec.Unschedulable,
			Conditions:              nodeHelper.GetNodeConditions(&item), // 获取节点事件
			Address:                 nodeHelper.GetNodeInternalIP(&item), // 获取节点内网 IP
			Status:                  nodeHelper.GetNodeStatus(&item),     // 获取节点状态
			CreatedAt:               item.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
		}
		data = append(data, node)
	}
	return data, nil
}

func (r *nodes) Get(ctx context.Context, id uint, name string) (*corev1.Node, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 node get failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *nodes) GetDetail(ctx context.Context, id uint, name string) (*Node, error) {
	var (
		data Node
		pod  Pod
	)

	nodeHelper := helper.NewNodeHelper()
	podHelper := helper.NewPodHelper()

	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	// 获取 node 信息
	node, err := client.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return &data, fmt.Errorf("kubernetes CoreV1 nodes get failed. err: %v", err)
	}

	data = Node{
		Name:                    node.Name,
		Labels:                  node.Labels,
		KubeletVersion:          node.Status.NodeInfo.KubeletVersion,
		KubeProxyVersion:        node.Status.NodeInfo.KubeProxyVersion,
		KernelVersion:           node.Status.NodeInfo.KernelVersion,
		ContainerRuntimeVersion: node.Status.NodeInfo.ContainerRuntimeVersion,
		OSSystem:                node.Status.NodeInfo.OperatingSystem,
		OSImage:                 node.Status.NodeInfo.OSImage,
		Architecture:            node.Status.NodeInfo.Architecture,
		MachineID:               node.Status.NodeInfo.MachineID,
		SystemUUID:              node.Status.NodeInfo.SystemUUID,
		BootID:                  node.Status.NodeInfo.BootID,
		PodCIDR:                 node.Spec.PodCIDR,
		Taints:                  node.Spec.Taints,
		Unschedulable:           &node.Spec.Unschedulable,
		Conditions:              nodeHelper.GetNodeConditions(node), // 获取节点事件
		Address:                 nodeHelper.GetNodeInternalIP(node), // 获取节点内网 IP
		Status:                  nodeHelper.GetNodeStatus(node),     // 获取节点状态
		CreatedAt:               node.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
	}
	//  获取节点 IP
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			data.Address = addr.Address
			break // 每个节点通常只有一个 InternalIP，所以找到后可以 break
		}
	}
	// 获取节点状态
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			data.Status = condition.Status
			break
		}
	}
	//  获取节点 IP
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			data.Address = addr.Address
			break // 每个节点通常只有一个 InternalIP，所以找到后可以 break
		}
	}
	// 获取节点状态
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			data.Status = condition.Status
			break
		}
	}

	// 获取本节点上的 pods 信息
	nodeName := data.Name
	pods, err := client.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		return &data, fmt.Errorf("kubernetes CoreV1 pods list failed. err: %v", err)
	}
	for _, item := range pods.Items {
		pod.Namespace = item.Namespace
		pod.Name = item.Name
		pod.PodIP = item.Status.PodIP
		pod.HostIP = item.Status.HostIP
		pod.NodeName = item.Spec.NodeName
		pod.Status = podHelper.GetPodStatus(item)             // Pod 运行状态
		pod.RestartCount = podHelper.GetMaxRestartCount(item) // Pod 重启次数（获取重启次数最多的容器)
		pod.CreatedAt = item.CreationTimestamp.Time.Format("2006-01-02 15:04:05")
		pod.Labels = item.Labels
		data.Pods = append(data.Pods, pod)
	}
	return &data, nil
}

func (r *nodes) Cordon(ctx context.Context, id uint, name string) (*corev1.Node, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	// 获取节点信息
	node, err := client.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 node get failed. err: %v", err)
	}
	node.Spec.Unschedulable = !node.Spec.Unschedulable // 调度状态(取反)
	response, err := client.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 node unschedulable update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}

func (r *nodes) Update(ctx context.Context, id uint, node *corev1.Node) (*corev1.Node, error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}
	response, err := client.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 node update failed. err: %v", err)
	}
	r.gvkFactory.Complete(response)
	return response, nil
}
