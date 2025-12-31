package helper

import (
	corev1 "k8s.io/api/core/v1"
)

// NodeHelper 节点相关工具函数
type NodeHelper struct{}

func NewNodeHelper() *NodeHelper {
	return &NodeHelper{}
}

// GetNodeInternalIP 获取节点内部 IP
func (h *NodeHelper) GetNodeInternalIP(node *corev1.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			return addr.Address
		}
	}
	return ""
}

// GetNodeExternalIP 获取节点外部 IP
func (h *NodeHelper) GetNodeExternalIP(node *corev1.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeExternalIP {
			return addr.Address
		}
	}
	return ""
}

// GetNodeStatus 获取节点状态
func (h *NodeHelper) GetNodeStatus(node *corev1.Node) corev1.ConditionStatus {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status
		}
	}
	return corev1.ConditionUnknown
}

// GetNodeStatusString 获取节点状态字符串
func (h *NodeHelper) GetNodeStatusString(node *corev1.Node) string {
	status := h.GetNodeStatus(node)
	switch status {
	case corev1.ConditionTrue:
		return "True"
	case corev1.ConditionFalse:
		return "False"
	default:
		return "Unknown"
	}
}

// NodeCondition 节点条件信息
type NodeCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	Reason             string `json:"reason,omitempty"`
	Message            string `json:"message,omitempty"`
	LastHeartbeatTime  string `json:"lastHeartbeatTime"`
	LastTransitionTime string `json:"lastTransitionTime"`
}

// GetNodeConditions 获取所有节点条件
func (h *NodeHelper) GetNodeConditions(node *corev1.Node) []NodeCondition {
	var conditions []NodeCondition
	for _, cond := range node.Status.Conditions {
		conditions = append(conditions, NodeCondition{
			Type:               string(cond.Type),
			Status:             string(cond.Status),
			Reason:             cond.Reason,
			Message:            cond.Message,
			LastHeartbeatTime:  cond.LastHeartbeatTime.Time.Format("2006-01-02 15:04:05"),
			LastTransitionTime: cond.LastTransitionTime.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return conditions
}

// GetNodeConditionByType 根据类型获取节点条件
func (h *NodeHelper) GetNodeConditionByType(node *corev1.Node, conditionType corev1.NodeConditionType) *NodeCondition {
	for _, cond := range node.Status.Conditions {
		if cond.Type == conditionType {
			return &NodeCondition{
				Type:               string(cond.Type),
				Status:             string(cond.Status),
				Reason:             cond.Reason,
				Message:            cond.Message,
				LastHeartbeatTime:  cond.LastHeartbeatTime.Time.Format("2006-01-02 15:04:05"),
				LastTransitionTime: cond.LastTransitionTime.Time.Format("2006-01-02 15:04:05"),
			}
		}
	}
	return nil
}
