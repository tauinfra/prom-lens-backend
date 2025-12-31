package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/repository"
)

// EventService 定义接口
type EventService interface {
	List(ctx context.Context, id int, ns, kind, name string) (events []repository.Event, err error)
}

type eventService struct {
	event repository.EventRepository
}

func NewEventService(event repository.EventRepository) EventService {
	return &eventService{event: event}
}

// List 列表
func (s *eventService) List(ctx context.Context, id int, ns, kind, name string) (events []repository.Event, err error) {
	return s.event.List(ctx, id, ns, kind, name)
}
