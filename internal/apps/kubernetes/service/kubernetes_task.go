package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

// TaskManager 定义接口
type TaskManager interface {
	List(ctx context.Context, id uint, ns string) ([]dto.Task, error)
	Get(ctx context.Context, id uint, ns, name string) (*v1.Task, error)
	Create(ctx context.Context, id uint, ns string, body *v1.Task) (*v1.Task, error)
	Update(ctx context.Context, id uint, ns string, body *v1.Task) (*v1.Task, error)
	Delete(ctx context.Context, id uint, ns string, name string) error
}

type taskManager struct {
	task repository.TaskRepository
}

func NewTaskManager(task repository.TaskRepository) TaskManager {
	return &taskManager{task: task}
}

// List 列表
func (s *taskManager) List(ctx context.Context, id uint, ns string) ([]dto.Task, error) {
	items, err := s.task.List(ctx, id, ns)
	if err != nil {
		return nil, err
	}
	return dto.ToTaskDTOs(items), nil
}

// Get 查询
func (s *taskManager) Get(ctx context.Context, id uint, ns, name string) (*v1.Task, error) {
	return s.task.Get(ctx, id, ns, name)
}

// Create 创建
func (s *taskManager) Create(ctx context.Context, id uint, ns string, body *v1.Task) (*v1.Task, error) {
	return s.task.Create(ctx, id, ns, body)
}

// Update 更新
func (s *taskManager) Update(ctx context.Context, id uint, ns string, body *v1.Task) (*v1.Task, error) {
	return s.task.Update(ctx, id, ns, body)
}

// Delete 删除
func (s *taskManager) Delete(ctx context.Context, id uint, ns, name string) (err error) {
	return s.task.Delete(ctx, id, ns, name)
}
