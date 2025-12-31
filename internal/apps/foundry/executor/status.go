package executor

import (
	"fmt"
	"strings"
	"time"
	"valyria-backend/internal/apps/foundry/model"
	"valyria-backend/internal/core/logger"
)

func (p *PipelineExecutor) updateReleaseStatus(status ReleaseStatus) {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	switch status {
	case StatusRunning:
		if p.release.StartedAt != nil {
			return
		}
		err := p.db.
			Where("task_id = ?", p.release.TaskID).
			Select("Status", "StartedAt").
			Updates(model.Release{
				Status:    string(StatusRunning),
				StartedAt: &now,
			}).Error

		if err != nil {
			logger.Errorf("Failed to update release %v to running: %v", p.release.TaskID, err)
			return
		}
		p.status = StatusRunning
		p.release.StartedAt = &now

	case StatusSuccess, StatusFailed:
		if p.release.FinishedAt != nil {
			return
		}
		err := p.db.
			Where("task_id = ?", p.release.TaskID).
			Select("Status", "FinishedAt").
			Updates(model.Release{
				Status:     string(status),
				FinishedAt: &now,
			}).Error

		if err != nil {
			logger.Errorf("Failed to update release %v to %s: %v", p.release.TaskID, status, err)
			return
		}
		p.status = status
		p.release.FinishedAt = &now
	}
}

func (p *PipelineExecutor) finishRelease(status ReleaseStatus, phase ReleasePhase) {
	p.updateReleaseStatus(status)
	p.logReleaseEvent(p.release.TaskID, status, phase)
}

func (p *PipelineExecutor) logReleaseEvent(taskID string, status ReleaseStatus, phase ReleasePhase) {
	now := time.Now().Format("2006-01-02 15:04:05")
	sep := strings.Repeat("=", 10)
	msg := fmt.Sprintf("%s Release TaskID: %s %s | Status: %s | Time: %s %s",
		sep, taskID, phase, status, now, sep,
	)
	// 写入 Redis
	if err := p.log.AppendLog(msg); err != nil {
		logger.Errorf("Failed to write log: %v", err)
	}
	// 本地控制台输出
	logger.Infof(msg)
}
