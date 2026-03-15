package executor

import "valyria-backend/internal/core/config"

func getPromRuleConfig() (string, string) {
	cfg := config.GetConfig()
	namespace := "monitoring"
	configmap := "prometheus-rules"
	if cfg == nil {
		return namespace, configmap
	}
	if cfg.Prometheus.Rule.Namespace != "" {
		namespace = cfg.Prometheus.Rule.Namespace
	}
	if cfg.Prometheus.Rule.ConfigMap != "" {
		configmap = cfg.Prometheus.Rule.ConfigMap
	}
	return namespace, configmap
}

func getPromTargetConfig() (string, string) {
	cfg := config.GetConfig()
	namespace := "monitoring"
	configmap := "prometheus-targets"
	if cfg == nil {
		return namespace, configmap
	}
	if cfg.Prometheus.Target.Namespace != "" {
		namespace = cfg.Prometheus.Target.Namespace
	}
	if cfg.Prometheus.Target.ConfigMap != "" {
		configmap = cfg.Prometheus.Target.ConfigMap
	}
	return namespace, configmap
}
