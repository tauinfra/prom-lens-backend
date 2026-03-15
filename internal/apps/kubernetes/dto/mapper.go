package dto

import (
	"valyria-backend/internal/apps/kubernetes/repository"

	networkingv1 "k8s.io/api/networking/v1"
)

func ToConfigMapDTOs(items []repository.ConfigMap) []ConfigMap {
	if len(items) == 0 {
		return nil
	}
	out := make([]ConfigMap, len(items))
	for i, item := range items {
		out[i] = ConfigMap{
			Name:      item.Name,
			Data:      item.Data,
			CreatedAt: item.CreatedAt,
		}
	}
	return out
}

func ToHpaDTOs(items []repository.HPA) []HPA {
	if len(items) == 0 {
		return nil
	}
	out := make([]HPA, len(items))
	for i, item := range items {
		out[i] = ToHpaDTO(item)
	}
	return out
}

func ToHpaDTO(item repository.HPA) HPA {
	return HPA{
		Namespace:        item.Namespace,
		Name:             item.Name,
		ScaleTargetKind:  item.ScaleTargetKind,
		ScaleTargetName:  item.ScaleTargetName,
		MinReplicas:      item.MinReplicas,
		MaxReplicas:      item.MaxReplicas,
		CurrentReplicas:  item.CurrentReplicas,
		DesiredReplicas:  item.DesiredReplicas,
		CreatedAt:        item.CreatedAt,
	}
}

func ToEventDTOs(items []repository.Event) []Event {
	if len(items) == 0 {
		return nil
	}
	out := make([]Event, len(items))
	for i, item := range items {
		out[i] = Event{
			Type:     item.Type,
			Reason:   item.Reason,
			Object:   item.Object,
			Source:   item.Source,
			Message:  item.Message,
			CreateAt: item.CreateAt,
		}
	}
	return out
}

func ToNamespaceDTOs(items []repository.Namespace) []Namespace {
	if len(items) == 0 {
		return nil
	}
	out := make([]Namespace, len(items))
	for i, item := range items {
		out[i] = Namespace{
			Name:      item.Name,
			Status:    item.Status,
			Labels:    item.Labels,
			CreatedAt: item.CreatedAt,
		}
	}
	return out
}

func ToNodeDTOs(items []repository.Node) []Node {
	if len(items) == 0 {
		return nil
	}
	out := make([]Node, len(items))
	for i, item := range items {
		out[i] = ToNodeDTO(item)
	}
	return out
}

func ToNodeDTO(item repository.Node) Node {
	conditions := make([]NodeCondition, 0, len(item.Conditions))
	for _, c := range item.Conditions {
		conditions = append(conditions, NodeCondition{
			Type:               c.Type,
			Status:             c.Status,
			Reason:             c.Reason,
			Message:            c.Message,
			LastHeartbeatTime:  c.LastHeartbeatTime,
			LastTransitionTime: c.LastTransitionTime,
		})
	}
	return Node{
		Name:                    item.Name,
		Address:                 item.Address,
		OSImage:                 item.OSImage,
		OSSystem:                item.OSSystem,
		Architecture:            item.Architecture,
		KubeletVersion:          item.KubeletVersion,
		KubeProxyVersion:        item.KubeProxyVersion,
		KernelVersion:           item.KernelVersion,
		ContainerRuntimeVersion: item.ContainerRuntimeVersion,
		PodCIDR:                 item.PodCIDR,
		Labels:                  item.Labels,
		Taints:                  item.Taints,
		Status:                  item.Status,
		MachineID:               item.MachineID,
		SystemUUID:              item.SystemUUID,
		BootID:                  item.BootID,
		Unschedulable:           item.Unschedulable,
		Conditions:              conditions,
		Pods:                    ToPodDTOs(item.Pods),
		CreatedAt:               item.CreatedAt,
	}
}

func ToPodDTOs(items []repository.Pod) []Pod {
	if len(items) == 0 {
		return nil
	}
	out := make([]Pod, len(items))
	for i, item := range items {
		out[i] = ToPodDTO(item)
	}
	return out
}

func ToPodDTO(item repository.Pod) Pod {
	containers := make([]Container, 0, len(item.Containers))
	for _, c := range item.Containers {
		containers = append(containers, Container{
			Name:         c.Name,
			Image:        c.Image,
			ImageID:      c.ImageID,
			Ready:        c.Ready,
			State:        c.State,
			RestartCount: c.RestartCount,
			StartedAt:    c.StartedAt,
		})
	}
	return Pod{
		Namespace:    item.Namespace,
		Name:         item.Name,
		PodIP:        item.PodIP,
		HostIP:       item.HostIP,
		NodeName:     item.NodeName,
		Status:       item.Status,
		RestartCount: item.RestartCount,
		Ready:        item.Ready,
		Containers:   containers,
		Conditions:   item.Conditions,
		Labels:       item.Labels,
		CreatedAt:    item.CreatedAt,
	}
}

func ToPersistentVolumeDTOs(items []repository.PersistentVolume) []PersistentVolume {
	if len(items) == 0 {
		return nil
	}
	out := make([]PersistentVolume, len(items))
	for i, item := range items {
		out[i] = PersistentVolume{
			Name:          item.Name,
			Capacity:      item.Capacity,
			AccessModes:   item.AccessModes,
			ReclaimPolicy: item.ReclaimPolicy,
			Status:        item.Status,
			Claim:         item.Claim,
			StorageClass:  item.StorageClass,
			Reason:        item.Reason,
			CreatedAt:     item.CreatedAt,
		}
	}
	return out
}

func ToPersistentVolumeClaimDTOs(items []repository.PersistentVolumeClaim) []PersistentVolumeClaim {
	if len(items) == 0 {
		return nil
	}
	out := make([]PersistentVolumeClaim, len(items))
	for i, item := range items {
		out[i] = PersistentVolumeClaim{
			Name:         item.Name,
			Capacity:     item.Capacity,
			AccessModes:  item.AccessModes,
			Status:       item.Status,
			StorageClass: item.StorageClass,
			Volume:       item.Volume,
			CreatedAt:    item.CreatedAt,
		}
	}
	return out
}

func ToReplicaSetDTOs(items []repository.ReplicaSet) []ReplicaSet {
	if len(items) == 0 {
		return nil
	}
	out := make([]ReplicaSet, len(items))
	for i, item := range items {
		out[i] = ReplicaSet{
			Namespace: item.Namespace,
			Name:      item.Name,
			Image:     item.Image,
			Revision:  item.Revision,
			CreatedAt: item.CreatedAt,
		}
	}
	return out
}

func ToSecretDTOs(items []repository.Secret) []Secret {
	if len(items) == 0 {
		return nil
	}
	out := make([]Secret, len(items))
	for i, item := range items {
		out[i] = Secret{
			Name:      item.Name,
			Type:      item.Type,
			Data:      item.Data,
			CreatedAt: item.CreatedAt,
		}
	}
	return out
}

func ToRoleDTOs(items []repository.Role) []Role {
	if len(items) == 0 {
		return nil
	}
	out := make([]Role, len(items))
	for i, item := range items {
		out[i] = Role{
			Name:      item.Name,
			Rules:     item.Rules,
			CreatedAt: item.CreatedAt,
		}
	}
	return out
}

func ToClusterRoleDTOs(items []repository.ClusterRole) []ClusterRole {
	if len(items) == 0 {
		return nil
	}
	out := make([]ClusterRole, len(items))
	for i, item := range items {
		out[i] = ClusterRole{
			Name:      item.Name,
			Rules:     item.Rules,
			CreatedAt: item.CreatedAt,
		}
	}
	return out
}

func ToRoleBindingDTOs(items []repository.RoleBinding) []RoleBinding {
	if len(items) == 0 {
		return nil
	}
	out := make([]RoleBinding, len(items))
	for i, item := range items {
		out[i] = RoleBinding{
			Name:      item.Name,
			RoleRef:   item.RoleRef,
			Subjects:  item.Subjects,
			CreatedAt: item.CreatedAt,
		}
	}
	return out
}

func ToClusterRoleBindingDTOs(items []repository.ClusterRoleBinding) []ClusterRoleBinding {
	if len(items) == 0 {
		return nil
	}
	out := make([]ClusterRoleBinding, len(items))
	for i, item := range items {
		out[i] = ClusterRoleBinding{
			Name:      item.Name,
			RoleRef:   item.RoleRef,
			Subjects:  item.Subjects,
			CreatedAt: item.CreatedAt,
		}
	}
	return out
}

func ToServiceAccountDTOs(items []repository.ServiceAccount) []ServiceAccount {
	if len(items) == 0 {
		return nil
	}
	out := make([]ServiceAccount, len(items))
	for i, item := range items {
		out[i] = ServiceAccount{
			Name:             item.Name,
			Secrets:          item.Secrets,
			ImagePullSecrets: item.ImagePullSecrets,
			CreatedAt:        item.CreatedAt,
		}
	}
	return out
}

func ToIngressDTOs(items []repository.Ingress) []Ingress {
	if len(items) == 0 {
		return nil
	}
	out := make([]Ingress, len(items))
	for i, item := range items {
		out[i] = Ingress{
			Name:      item.Name,
			Hosts:     item.Hosts,
			Paths:     item.Paths,
			TLS:       item.TLS,
			Endpoint:  item.Endpoint,
			ClassName: item.ClassName,
			CreatedAt: item.CreatedAt,
		}
	}
	return out
}

func ToIngressDetailDTO(item *networkingv1.Ingress) IngressDetail {
	out := IngressDetail{
		Name:        item.Name,
		Namespace:   item.Namespace,
		Annotations: filterMap(item.Annotations, "kubectl.kubernetes.io/last-applied-configuration"),
		Labels:      item.Labels,
		CreatedAt:   item.CreationTimestamp.UTC().Format("2006-01-02T15:04:05Z"),
	}
	if item.Spec.IngressClassName != nil {
		out.ClassName = *item.Spec.IngressClassName
	}

	rules := make([]IngressRule, 0, len(item.Spec.Rules))
	for _, rule := range item.Spec.Rules {
		r := IngressRule{
			Host: rule.Host,
		}
		if r.Host == "" {
			r.Host = "*"
		}
		if rule.HTTP != nil {
			paths := make([]IngressRulePath, 0, len(rule.HTTP.Paths))
			for _, p := range rule.HTTP.Paths {
				pathType := ""
				if p.PathType != nil {
					pathType = string(*p.PathType)
				}
				path := p.Path
				if path == "" {
					path = "/"
				}
				port := 0
				if p.Backend.Service != nil {
					port = int(p.Backend.Service.Port.Number)
				}
				service := ""
				if p.Backend.Service != nil {
					service = p.Backend.Service.Name
				}
				paths = append(paths, IngressRulePath{
					Path:     path,
					PathType: pathType,
					Service:  service,
					Port:     port,
				})
			}
			r.Paths = paths
		}
		rules = append(rules, r)
	}
	out.Rules = rules

	return out
}

func filterMap(src map[string]string, excludes ...string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	excluded := make(map[string]struct{}, len(excludes))
	for _, k := range excludes {
		excluded[k] = struct{}{}
	}
	out := make(map[string]string, len(src))
	for k, v := range src {
		if _, skip := excluded[k]; skip {
			continue
		}
		out[k] = v
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func ToIngressClassDTOs(items []repository.IngressClass) []IngressClass {
	if len(items) == 0 {
		return nil
	}
	out := make([]IngressClass, len(items))
	for i, item := range items {
		out[i] = IngressClass{
			Name:       item.Name,
			Controller: item.Controller,
			Parameters: item.Parameters,
			CreatedAt:  item.CreatedAt,
		}
	}
	return out
}

func ToStorageClassDTOs(items []repository.StorageClass) []StorageClass {
	if len(items) == 0 {
		return nil
	}
	out := make([]StorageClass, len(items))
	for i, item := range items {
		out[i] = StorageClass{
			Name:                 item.Name,
			Provisioner:          item.Provisioner,
			ReclaimPolicy:        item.ReclaimPolicy,
			VolumeBindingMode:    item.VolumeBindingMode,
			AllowVolumeExpansion: item.AllowVolumeExpansion,
			CreatedAt:            item.CreatedAt,
		}
	}
	return out
}

func ToStatefulSetDTOs(items []repository.StatefulSet) []StatefulSet {
	if len(items) == 0 {
		return nil
	}
	out := make([]StatefulSet, len(items))
	for i, item := range items {
		out[i] = ToStatefulSetDTO(item)
	}
	return out
}

func ToStatefulSetDTO(item repository.StatefulSet) StatefulSet {
	return StatefulSet{
		Namespace:         item.Namespace,
		Name:              item.Name,
		Replicas:          item.Replicas,
		ReadyReplicas:     item.ReadyReplicas,
		CurrentReplicas:   item.CurrentReplicas,
		AvailableReplicas: item.AvailableReplicas,
		UpdatedReplicas:   item.UpdatedReplicas,
		MatchLabels:       item.MatchLabels,
		Strategy:          item.Strategy,
		CreatedAt:         item.CreatedAt,
	}
}

func ToDaemonSetDTOs(items []repository.DaemonSet) []DaemonSet {
	if len(items) == 0 {
		return nil
	}
	out := make([]DaemonSet, len(items))
	for i, item := range items {
		out[i] = ToDaemonSetDTO(item)
	}
	return out
}

func ToDaemonSetDTO(item repository.DaemonSet) DaemonSet {
	return DaemonSet{
		Namespace:              item.Namespace,
		Name:                   item.Name,
		CurrentNumberScheduled: item.CurrentNumberScheduled,
		DesiredNumberScheduled: item.DesiredNumberScheduled,
		NumberReady:            item.NumberReady,
		NumberAvailable:        item.NumberAvailable,
		UpdatedNumberScheduled: item.UpdatedNumberScheduled,
		MatchLabels:            item.MatchLabels,
		Strategy:               item.Strategy,
		CreatedAt:              item.CreatedAt,
	}
}

func ToDeploymentDTOs(items []repository.Deployment) []Deployment {
	if len(items) == 0 {
		return nil
	}
	out := make([]Deployment, len(items))
	for i, item := range items {
		out[i] = ToDeploymentDTO(item)
	}
	return out
}

func ToDeploymentDTO(item repository.Deployment) Deployment {
	return Deployment{
		Namespace:         item.Namespace,
		Name:              item.Name,
		Replicas:          item.Replicas,
		ReadyReplicas:     item.ReadyReplicas,
		AvailableReplicas: item.AvailableReplicas,
		UpdatedReplicas:   item.UpdatedReplicas,
		Paused:            item.Paused,
		MatchLabels:       item.MatchLabels,
		Strategy:          item.Strategy,
		Conditions:        item.Conditions,
		CreatedAt:         item.CreatedAt,
		UpdatedAt:         item.UpdatedAt,
	}
}

func ToServiceDTOs(items []repository.Service) []Service {
	if len(items) == 0 {
		return nil
	}
	out := make([]Service, len(items))
	for i, item := range items {
		out[i] = Service{
			Name:      item.Name,
			Type:      item.Type,
			ClusterIP: item.ClusterIP,
			Ports:     item.Ports,
			CreatedAt: item.CreatedAt,
		}
	}
	return out
}

func ToPipelineDTOs(items []repository.Pipeline) []Pipeline {
	if len(items) == 0 {
		return nil
	}
	out := make([]Pipeline, len(items))
	for i, item := range items {
		out[i] = Pipeline{
			Name:      item.Name,
			Namespace: item.Namespace,
			CreatedAt: item.CreatedAt,
		}
	}
	return out
}

func ToPipelineRunDTOs(items []repository.PipelineRun) []PipelineRun {
	if len(items) == 0 {
		return nil
	}
	out := make([]PipelineRun, len(items))
	for i, item := range items {
		out[i] = PipelineRun{
			Name:      item.Name,
			Namespace: item.Namespace,
			CreatedAt: item.CreatedAt,
		}
	}
	return out
}

func ToTaskDTOs(items []repository.Task) []Task {
	if len(items) == 0 {
		return nil
	}
	out := make([]Task, len(items))
	for i, item := range items {
		out[i] = Task{
			Name:      item.Name,
			Namespace: item.Namespace,
			CreatedAt: item.CreatedAt,
		}
	}
	return out
}

func ToTaskRunDTOs(items []repository.TaskRun) []TaskRun {
	if len(items) == 0 {
		return nil
	}
	out := make([]TaskRun, len(items))
	for i, item := range items {
		out[i] = TaskRun{
			Name:      item.Name,
			Namespace: item.Namespace,
			Status:    item.Status,
			Reason:    item.Reason,
			StartAt:   item.StartAt,
			EndAt:     item.EndAt,
			CreatedAt: item.CreatedAt,
		}
	}
	return out
}
