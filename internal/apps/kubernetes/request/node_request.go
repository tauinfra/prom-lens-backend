package request

import corev1 "k8s.io/api/core/v1"

type UpdateNodeRequest struct {
	Labels map[string]string `json:"labels,omitempty"`
	Taints []corev1.Taint    `json:"taints,omitempty"`
}
