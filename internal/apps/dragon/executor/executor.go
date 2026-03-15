package executor

import (
	"context"
	"fmt"
	"sync"
	"time"
	"valyria-backend/internal/apps/dragon/adapter"
	"valyria-backend/internal/apps/dragon/model"
	"valyria-backend/internal/core/config"
	"valyria-backend/internal/core/logger"

	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PipelineExecutor struct {
	ctx            context.Context
	tektonClient   *tektonclient.Clientset
	client         *kubernetes.Clientset
	db             *gorm.DB
	release        model.Release
	taskOrder      []string
	taskOrderMutex sync.RWMutex
	status         ReleaseStatus
	statusMutex    sync.Mutex
	logWriter      *logWriter
	processedPods  sync.Map // key = pod/container
	cfg            *config.Config
	repoURL        string
}

func Run(ctx context.Context, taskID, repoURL string, db *gorm.DB, cfg *config.Config) (*PipelineExecutor, error) {
	var release model.Release
	if err := db.Where("task_id = ?", taskID).
		Preload("Pipeline.Environment.Project").
		First(&release).Error; err != nil {
		return nil, err
	}
	// k8s client
	tektonClient, client, err := adapter.NewTektonClients()
	if err != nil {
		return nil, err
	}
	h := &PipelineExecutor{
		ctx:          ctx,
		repoURL:      repoURL,
		tektonClient: tektonClient,
		client:       client,
		db:           db,
		release:      release,
		status:       StatusPending,
		logWriter:    newLogWriter(release.TaskID),
		cfg:          cfg,
	}
	go h.Create()
	return h, nil
}

func (p *PipelineExecutor) Create() {
	defer p.cleanup()
	pr := p.pipelineRunTemplate()
	// 更新 Release 状态
	p.finishRelease(StatusRunning, PhaseStarting)
	// 1. 创建 PipelineRun
	_, err := p.tektonClient.TektonV1().
		PipelineRuns(p.cfg.K8s.Tekton.Namespace).
		Create(p.ctx, pr, metav1.CreateOptions{})
	if err != nil {
		p.finishRelease(StatusFailed, PhaseCompleted)
		logger.Errorf("tekton pipelineRun '%v' creation failed. error: %v", pr.Name, err)
		logMsg := fmt.Sprintf("[%s]  %s\n", time.Now().Format("2006-01-02 15:04:05"), err)
		if err = p.logWriter.AppendLog(logMsg); err != nil {
			logger.Errorf("Failed to write log: %v", err)
		}
		return
	}
	// 2. 检查 PipelineRun 运行状态
	if err = p.waitPipelineInitialized(pr.Name); err != nil {
		p.finishRelease(StatusFailed, PhaseCompleted)
		logger.Errorf("tekton pipelineRun '%v' failed to start. error: %v", pr.Name, err)
		logMsg := fmt.Sprintf("[%s]  %s\n", time.Now().Format("2006-01-02 15:04:05"), err)
		if err = p.logWriter.AppendLog(logMsg); err != nil {
			logger.Errorf("Failed to write log: %v", err)
		}
		return
	}
	logger.Infof("tekton pipelineRun '%s' created successfully.", pr.Name)
	// 3. 观察 TaskRun
	if err = p.waitTaskRuns(pr.Name); err != nil {
		p.finishRelease(StatusFailed, PhaseCompleted)
		logMsg := fmt.Sprintf("[%s]  %s\n", time.Now().Format("2006-01-02 15:04:05"), err)
		if err = p.logWriter.AppendLog(logMsg); err != nil {
			logger.Errorf("Failed to write log: %v", err)
		}
		return
	}
	logger.Infof("tekton PipelineRun '%s' completed successfully.", pr.Name)
	p.finishRelease(StatusSuccess, PhaseCompleted)
}

// cleanup 清理执行器资源
func (p *PipelineExecutor) cleanup() {
	// 清理 processedPods，防止内存泄漏
	p.processedPods.Range(func(key, value interface{}) bool {
		p.processedPods.Delete(key)
		return true
	})
}
