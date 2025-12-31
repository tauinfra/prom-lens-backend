package service

import (
	"context"
	"io"
	"os"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/core/logger"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/remotecommand"
)

type PodService interface {
	List(ctx context.Context, id int, ns, labelSelector string) ([]repository.Pod, error)
	Get(ctx context.Context, id int, ns, name string) (pod *corev1.Pod, err error)
	GetDetail(ctx context.Context, id int, ns, name string) (repository.Pod, error)
	GetLogs(ctx context.Context, id int, ns, name, container string, follow bool, tailLines, sinceSeconds *int64) (io.ReadCloser, context.CancelFunc, error)
	Delete(ctx context.Context, id int, ns, name string) error
	Executor(ctx context.Context, id int, ns, name, container string) (remotecommand.Executor, error)
	DebugExecutor(ctx context.Context, id int, ns, name, container string) error
}

type podService struct {
	pod repository.PodRepository
}

func NewPodService(pod repository.PodRepository) PodService {
	return &podService{pod: pod}
}

// List 列表
func (s *podService) List(ctx context.Context, id int, ns, labelSelector string) ([]repository.Pod, error) {
	return s.pod.List(ctx, id, ns, labelSelector)
}

// Get 查询
func (s *podService) Get(ctx context.Context, id int, ns, name string) (pod *corev1.Pod, err error) {
	return s.pod.Get(ctx, id, ns, name)
}

// GetDetail 详情
func (s *podService) GetDetail(ctx context.Context, id int, ns, name string) (repository.Pod, error) {
	return s.pod.GetDetail(ctx, id, ns, name)
}

// GetLogs 日志
func (s *podService) GetLogs(ctx context.Context, id int, ns, name, container string, follow bool, tailLines, sinceSeconds *int64) (io.ReadCloser, context.CancelFunc, error) {
	return s.pod.GetLogs(ctx, id, ns, name, container, follow, tailLines, sinceSeconds)
}

// Delete 删除
func (s *podService) Delete(ctx context.Context, id int, ns, name string) error {
	return s.pod.Delete(ctx, id, ns, name)
}

// Executor 登录
func (s *podService) Executor(ctx context.Context, id int, ns, name, container string) (remotecommand.Executor, error) {
	return s.pod.Executor(ctx, id, ns, name, container)
}

func (s *podService) DebugExecutor(ctx context.Context, id int, ns, name, container string) (err error) {
	var executor remotecommand.Executor
	// 执行远程命令
	executor, err = s.pod.Executor(ctx, id, ns, name, container)
	if err = executor.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	}); err != nil {
		logger.Errorf("executor failed, error: %v\n", err)
		return err
	}
	return nil
}
