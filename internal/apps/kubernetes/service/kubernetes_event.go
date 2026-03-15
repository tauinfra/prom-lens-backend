package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/repository"
)

// EventManager 定义接口
type EventManager interface {
	List(ctx context.Context, id uint, ns, kind, name string) (events []dto.Event, err error)
}

type eventManager struct {
	event repository.EventRepository
}

func NewEventManager(event repository.EventRepository) EventManager {
	return &eventManager{event: event}
}

// List 列表
func (s *eventManager) List(ctx context.Context, id uint, ns, kind, name string) (events []dto.Event, err error) {
	items, err := s.event.List(ctx, id, ns, kind, name)
	if err != nil {
		return nil, err
	}
	return dto.ToEventDTOs(items), nil
}
