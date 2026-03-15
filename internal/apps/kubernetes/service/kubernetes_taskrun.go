package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

// TaskRunManager 定义接口
type TaskRunManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.TaskRun, error)
	Get(ctx context.Context, id uint, ns, name string) (*v1.TaskRun, error)
	Create(ctx context.Context, id uint, ns string, body *v1.TaskRun) (*v1.TaskRun, error)
	Update(ctx context.Context, id uint, ns string, body *v1.TaskRun) (*v1.TaskRun, error)
	Delete(ctx context.Context, id uint, ns string, name string) error
}

type taskRunManager struct {
	taskRun repository.TaskRunRepository
}

func NewTaskRunManager(taskRun repository.TaskRunRepository) TaskRunManager {
	return &taskRunManager{taskRun: taskRun}
}

// List 列表
func (s *taskRunManager) List(ctx context.Context, id uint, ns string) ([]dto.TaskRun, error) {
	items, err := s.taskRun.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToTaskRunDTOs(items), nil
}

// Get 查询
func (s *taskRunManager) Get(ctx context.Context, id uint, ns, name string) (*v1.TaskRun, error) {
	return s.taskRun.Get(ctx, id, ns, name)
}

// Create 创建
func (s *taskRunManager) Create(ctx context.Context, id uint, ns string, body *v1.TaskRun) (*v1.TaskRun, error) {
	return s.taskRun.Create(ctx, id, ns, body)
}

// Update 更新
func (s *taskRunManager) Update(ctx context.Context, id uint, ns string, body *v1.TaskRun) (*v1.TaskRun, error) {
	return s.taskRun.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *taskRunManager) Delete(ctx context.Context, id uint, ns, name string) (err error) {
	return s.taskRun.Delete(ctx, id, ns, name)
}
