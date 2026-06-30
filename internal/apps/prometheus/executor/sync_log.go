package executor

import "prom-lens-backend/internal/core/logger"

// logPromSyncFailed 记录 Prometheus ConfigMap 同步失败。
func logPromSyncFailed(resource string, groupID int, groupName, step string, err error) {
	logger.Errorf(
		"[prom-sync] failed resource=%s groupID=%d groupName=%q step=%s err=%v",
		resource, groupID, groupName, step, err,
	)
}

// logPromSyncSuccess 记录 Prometheus ConfigMap 同步成功。
func logPromSyncSuccess(resource, namespace, configMap, key string, payloadBytes int) {
	logger.Infof(
		"[prom-sync] success resource=%s namespace=%s configMap=%s key=%s bytes=%d",
		resource, namespace, configMap, key, payloadBytes,
	)
}
