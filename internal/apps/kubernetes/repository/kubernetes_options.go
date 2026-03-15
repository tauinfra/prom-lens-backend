package repository

import "valyria-backend/internal/core/config"

func timeoutSeconds() *int64 {
	cfg := config.GetConfig()
	timeout := int64(30) // 默认超时 30 秒
	if cfg != nil && cfg.K8s.Timeout > 0 {
		timeout = int64(cfg.K8s.Timeout)
	}
	return &timeout
}

func kubeClientLimits() (float32, int) {
	cfg := config.GetConfig()
	qps := float32(50) // 默认 QPS
	burst := 100       // 默认 Burst
	if cfg != nil && cfg.K8s.QPS > 0 {
		qps = cfg.K8s.QPS
	}
	if cfg != nil && cfg.K8s.Burst > 0 {
		burst = cfg.K8s.Burst
	}
	return qps, burst
}
