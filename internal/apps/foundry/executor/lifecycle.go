package executor

import (
	"context"
	"valyria-backend/internal/apps/foundry/model"
	"valyria-backend/internal/core/logger"

	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Run(taskID string, db *gorm.DB) (*PipelineExecutor, error) {
	var release model.Release
	if err := db.Where("task_id = ?", taskID).
		Preload("Pipeline.Environment.Project").
		Preload("Pipeline.Application").
		First(&release).Error; err != nil {
		return nil, err
	}

	tektonClient, kubeClient, err := newClients()
	if err != nil {
		return nil, err
	}

	h := &PipelineExecutor{
		tektonClient: tektonClient,
		kubeClient:   kubeClient,
		db:           db,
		release:      release,
		status:       StatusPending,
		log:          newLogWriter(release.TaskID),
	}
	go h.Create()
	return h, nil
}

func (t *PipelineExecutor) Create() {
	pr := t.pipelineRunTemplate()
	// 更新 Release 状态
	t.finishRelease(StatusRunning, PhaseStarting)
	// 1. 创建 PipelineRun
	_, err := t.tektonClient.TektonV1().
		PipelineRuns("default").
		Create(context.Background(), pr, metav1.CreateOptions{})
	if err != nil {
		logger.Errorf("tekton pipelineRun '%v' creation failed. error: %v", pr.Name, err)
		t.finishRelease(StatusFailed, PhaseCompleted)
		return
	}
	logger.Infof("tekton pipelineRun '%s' created successfully.", pr.Name)
	// 2. 检查 PipelineRun 运行状态
	if err = t.waitPipelineInitialized(pr.Name); err != nil {
		t.finishRelease(StatusFailed, PhaseCompleted)
		return
	}
	// 3. 观察 TaskRun
	if err := t.waitTaskRuns(pr.Name); err != nil {
		t.finishRelease(StatusFailed, PhaseCompleted)
		return
	}
	logger.Infof("tekton PipelineRun '%s' completed successfully.", pr.Name)
	t.finishRelease(StatusSuccess, PhaseCompleted)
}
