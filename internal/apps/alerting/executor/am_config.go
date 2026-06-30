package executor

import (
	"strings"

	"prom-lens-backend/internal/core/config"
)

func getAlertmanagerConfig() (namespace, configMap, configKey, baseURL, defaultReceiver string, enabled bool) {
	cfg := config.GetConfig()
	if cfg == nil {
		return "", "", "", "", "", false
	}
	baseURL = trimRightSlash(cfg.BaseURL)
	defaultReceiver = cfg.Alerting.Alertmanager.DefaultReceiver
	am := cfg.Alerting.Alertmanager
	if am.Namespace == "" || am.ConfigMap == "" || am.ConfigKey == "" || baseURL == "" {
		return am.Namespace, am.ConfigMap, am.ConfigKey, baseURL, defaultReceiver, false
	}
	return am.Namespace, am.ConfigMap, am.ConfigKey, baseURL, defaultReceiver, true
}

// IsAlertmanagerSyncConfigured 是否已配置 Alertmanager 同步。
func IsAlertmanagerSyncConfigured() bool {
	_, _, _, _, _, ok := getAlertmanagerConfig()
	return ok
}

func trimRightSlash(s string) string {
	s = strings.TrimSpace(s)
	return strings.TrimRight(s, "/")
}
