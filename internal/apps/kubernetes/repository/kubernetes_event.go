package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/pkg/k8s/helper"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EventRepository interface {
	List(ctx context.Context, id uint, ns, kind, name string) (events []Event, err error)
}

type eventRepository struct {
	cfgFactory *KubeConfigFactory
}

func NewEventRepository(kubeFactory *KubeConfigFactory) EventRepository {
	return &eventRepository{cfgFactory: kubeFactory}
}

type Event struct {
	Type     string `json:"type"`
	Reason   string `json:"reason"`
	Object   string `json:"object"`
	Source   string `json:"source"`
	Message  string `json:"message"`
	CreateAt string `json:"createAt"`
}

func (r *eventRepository) List(ctx context.Context, id uint, ns, kind, name string) (events []Event, err error) {
	client, err := r.cfgFactory.GetClientSet(ctx, id)
	if err != nil {
		return nil, err
	}

	eventHelper := helper.NewEventHelper()
	response, err := client.CoreV1().Events(ns).List(ctx, metav1.ListOptions{
		FieldSelector: eventHelper.BuildEventFieldSelector(kind, name),
	})
	if err != nil {
		return nil, fmt.Errorf("kubernetes CoreV1 events list failed. err: %v", err)
	}
	for _, item := range response.Items {
		events = append(events, Event{
			Type:     item.Type,
			Reason:   item.Reason,
			Object:   fmt.Sprintf("%v/%v", item.InvolvedObject.Kind, item.InvolvedObject.Name),
			Source:   item.Source.Component,
			Message:  item.Message,
			CreateAt: item.LastTimestamp.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return events, nil
}
