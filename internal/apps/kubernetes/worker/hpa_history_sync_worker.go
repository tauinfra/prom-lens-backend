package worker

import (
	"context"
	"time"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/core/logger"
)

// HpaHistorySyncWorker 定时从各集群拉取 HPA 相关 Event 并写入扩缩容历史表
type HpaHistorySyncWorker struct {
	Svc      service.HpaHistoryManager
	Interval time.Duration
}

// NewHpaHistorySyncWorker 创建 worker，interval 建议 1～2 分钟
func NewHpaHistorySyncWorker(svc service.HpaHistoryManager, interval time.Duration) *HpaHistorySyncWorker {
	if interval <= 0 {
		interval = 2 * time.Minute
	}
	return &HpaHistorySyncWorker{Svc: svc, Interval: interval}
}

// Run 阻塞运行，直到 ctx 取消
func (w *HpaHistorySyncWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	logger.Infof("HpaHistorySyncWorker started, interval=%v", w.Interval)
	for {
		select {
		case <-ctx.Done():
			logger.Info("HpaHistorySyncWorker stopped")
			return
		case <-ticker.C:
			if err := w.Svc.SyncFromClusters(ctx); err != nil {
				logger.Errorf("HpaHistorySyncWorker sync: %v", err)
			}
		}
	}
}
