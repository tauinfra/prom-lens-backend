package worker

import (
	"context"
	"time"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/core/logger"
)

// PermissionSyncWorker 后台轮询待同步的 K8s 权限并重试（幂等：仅 next_retry_at <= now 的记录）
type PermissionSyncWorker struct {
	Svc      service.PermissionManager
	Interval time.Duration
}

// NewPermissionSyncWorker 创建 worker，interval 建议 1～2 分钟
func NewPermissionSyncWorker(svc service.PermissionManager, interval time.Duration) *PermissionSyncWorker {
	if interval <= 0 {
		interval = 2 * time.Minute
	}
	return &PermissionSyncWorker{Svc: svc, Interval: interval}
}

// Run 阻塞运行，直到 ctx 取消
func (w *PermissionSyncWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	logger.Infof("PermissionSyncWorker started, interval=%v", w.Interval)
	for {
		select {
		case <-ctx.Done():
			logger.Info("PermissionSyncWorker stopped")
			return
		case <-ticker.C:
			if err := w.Svc.SyncPendingPermissions(ctx); err != nil {
				logger.Errorf("PermissionSyncWorker sync: %v", err)
			}
		}
	}
}
