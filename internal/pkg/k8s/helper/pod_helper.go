package helper

import (
	corev1 "k8s.io/api/core/v1"
)

// PodHelper Pod 相关工具函数
type PodHelper struct{}

type Container struct {
	Name         string `json:"name"`
	Image        string `json:"image"`
	ImageID      string `json:"imageID"`
	Ready        bool   `json:"ready"`
	State        string `json:"state"`
	RestartCount int32  `json:"restartCount"`
	ContainerID  string `json:"containerID"`
}

func NewPodHelper() *PodHelper {
	return &PodHelper{}
}

// GetPodStatus 获得 Pod 状态
func (h *PodHelper) GetPodStatus(pod corev1.Pod) string {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.State.Waiting != nil {
			return containerStatus.State.Waiting.Reason
		}
		if containerStatus.State.Terminated != nil {
			return containerStatus.State.Terminated.Reason
		}
	}
	return string(pod.Status.Phase)
}

// GetMaxRestartCount 获取 Pod 最大重启次数
func (h *PodHelper) GetMaxRestartCount(pod corev1.Pod) int {
	maxRestartCount := 0
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if int(containerStatus.RestartCount) > maxRestartCount {
			maxRestartCount = int(containerStatus.RestartCount)
		}
	}
	return maxRestartCount
}

// GetMaxReadyCount 获取容器运行数量
func (h *PodHelper) GetMaxReadyCount(pod corev1.Pod) int {
	maxReadyCount := 0
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Ready {
			maxReadyCount = +1
		}
	}
	return maxReadyCount
}

// GetContainerState 获取容器运行状态
func (h *PodHelper) GetContainerState(cs corev1.ContainerStatus) string {
	var state string
	switch {
	case cs.State.Running != nil:
		state = "Running"
	case cs.State.Waiting != nil:
		state = cs.State.Waiting.Reason
	case cs.State.Terminated != nil:
		state = cs.State.Terminated.Reason
	default:
		state = "unknown"
	}
	return state
}

// GetContainerStatuses 获取容器状态
func (h *PodHelper) GetContainerStatuses(pod corev1.Pod) []Container {
	var containers []Container
	for _, cs := range pod.Status.ContainerStatuses {
		container := Container{
			Name:         cs.Name,
			Ready:        cs.Ready,
			RestartCount: cs.RestartCount,
			Image:        cs.Image,
			ImageID:      cs.ImageID,
			ContainerID:  cs.ContainerID,
		}
		container.State = h.GetContainerState(cs)
		containers = append(containers, container)
	}
	return containers
}
