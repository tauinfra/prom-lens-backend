package service

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/request"
	"valyria-backend/internal/core/logger"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/remotecommand"
)

type PodManager interface {
	List(ctx context.Context, id uint, ns, labelSelector string) ([]dto.Pod, error)
	ListAll(ctx context.Context, id uint, labelSelector, fieldSelector string) ([]dto.Pod, error)
	Get(ctx context.Context, id uint, ns, name string) (pod *corev1.Pod, err error)
	GetDetail(ctx context.Context, id uint, ns, name string) (dto.Pod, error)
	GetLogs(ctx context.Context, id uint, ns, name, container string, follow bool, tailLines, sinceSeconds *int64) (io.ReadCloser, context.CancelFunc, error)
	Delete(ctx context.Context, id uint, ns, name string) error
	DeleteBatch(ctx context.Context, id uint, ns string, req *request.BatchDeletePodRequest) error
	Executor(ctx context.Context, id uint, ns, name, container string) (remotecommand.Executor, error)
	DebugExecutor(ctx context.Context, id uint, ns, name, container string) error
}

type podManager struct {
	pod repository.PodRepository
}

func NewPodManager(pod repository.PodRepository) PodManager {
	return &podManager{pod: pod}
}

// List 列表
func (s *podManager) List(ctx context.Context, id uint, ns, labelSelector string) ([]dto.Pod, error) {
	items, err := s.pod.List(ctx, id, ns, labelSelector)
	if err != nil {
		return nil, err
	}
	return dto.ToPodDTOs(items), nil
}

// ListAll 集群级别 Pod 列表，支持 fieldSelector、labelSelector
func (s *podManager) ListAll(ctx context.Context, id uint, labelSelector, fieldSelector string) ([]dto.Pod, error) {
	items, err := s.pod.ListAll(ctx, id, labelSelector, fieldSelector)
	if err != nil {
		return nil, err
	}
	return dto.ToPodDTOs(items), nil
}

// Get 查询
func (s *podManager) Get(ctx context.Context, id uint, ns, name string) (pod *corev1.Pod, err error) {
	return s.pod.Get(ctx, id, ns, name)
}

// GetDetail 详情
func (s *podManager) GetDetail(ctx context.Context, id uint, ns, name string) (dto.Pod, error) {
	item, err := s.pod.GetDetail(ctx, id, ns, name)
	if err != nil {
		return dto.Pod{}, err
	}
	return dto.ToPodDTO(item), nil
}

// GetLogs 日志
func (s *podManager) GetLogs(ctx context.Context, id uint, ns, name, container string, follow bool, tailLines, sinceSeconds *int64) (io.ReadCloser, context.CancelFunc, error) {
	return s.pod.GetLogs(ctx, id, ns, name, container, follow, tailLines, sinceSeconds)
}

// Delete 删除
func (s *podManager) Delete(ctx context.Context, id uint, ns, name string) error {
	return s.pod.Delete(ctx, id, ns, name)
}

// DeleteBatch 批量删除
func (s *podManager) DeleteBatch(ctx context.Context, id uint, ns string, req *request.BatchDeletePodRequest) error {
	if req == nil || len(req.Names) == 0 {
		return errors.New("pod names is required")
	}
	seen := make(map[string]struct{}, len(req.Names))
	for _, name := range req.Names {
		name = strings.TrimSpace(name)
		if name == "" {
			return errors.New("pod name is required")
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		if err := s.pod.Delete(ctx, id, ns, name); err != nil {
			return err
		}
	}
	return nil
}

// Executor 登录
func (s *podManager) Executor(ctx context.Context, id uint, ns, name, container string) (remotecommand.Executor, error) {
	return s.pod.Executor(ctx, id, ns, name, container)
}

func (s *podManager) DebugExecutor(ctx context.Context, id uint, ns, name, container string) (err error) {
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
