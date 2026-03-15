package service

import (
	"context"
	"valyria-backend/internal/apps/kubernetes/model"
	"valyria-backend/internal/apps/kubernetes/repository"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/core/logger"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const maxMessageLen = 500

// HpaHistoryManager HPA 扩缩容历史：查询 + 从集群 Event 同步
type HpaHistoryManager interface {
	ListHistory(ctx context.Context, clusterID uint, namespace, hpaName string, page, size int) ([]model.HPAScalingHistory, int64, error)
	SyncFromClusters(ctx context.Context) error
}

type hpaHistoryManager struct {
	historyRepo repository.HpaHistoryRepository
	cfgFactory  *repository.KubeConfigFactory
	clusterRepo repository.ClusterRepository
}

func NewHpaHistoryManager(historyRepo repository.HpaHistoryRepository, cfgFactory *repository.KubeConfigFactory, clusterRepo repository.ClusterRepository) HpaHistoryManager {
	return &hpaHistoryManager{
		historyRepo: historyRepo,
		cfgFactory:  cfgFactory,
		clusterRepo: clusterRepo,
	}
}

func (s *hpaHistoryManager) ListHistory(ctx context.Context, clusterID uint, namespace, hpaName string, page, size int) ([]model.HPAScalingHistory, int64, error) {
	return s.historyRepo.List(ctx, clusterID, namespace, hpaName, page, size)
}

func (s *hpaHistoryManager) SyncFromClusters(ctx context.Context) error {
	clusters, _, err := s.clusterRepo.List(ctx, pg.QueryParams{Page: 1, Size: 500})
	if err != nil {
		return err
	}
	for _, c := range clusters {
		if err := s.syncCluster(ctx, uint(c.ID)); err != nil {
			logger.Warnf("HPA history sync cluster %d: %v", c.ID, err)
		}
	}
	return nil
}

func (s *hpaHistoryManager) syncCluster(ctx context.Context, clusterID uint) error {
	client, err := s.cfgFactory.GetClientSet(ctx, clusterID)
	if err != nil {
		return err
	}
	nsList, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{Limit: 500})
	if err != nil {
		return err
	}
	for _, ns := range nsList.Items {
		namespace := ns.Name
		eventList, err := client.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{
			FieldSelector: "involvedObject.kind=HorizontalPodAutoscaler",
			Limit:         100,
		})
		if err != nil {
			logger.Warnf("HPA history list events cluster=%d ns=%s: %v", clusterID, namespace, err)
			continue
		}
		for i := range eventList.Items {
			ev := &eventList.Items[i]
			if err := s.ingestEvent(ctx, clusterID, ev); err != nil {
				logger.Warnf("HPA history ingest event %s: %v", ev.UID, err)
			}
		}
	}
	return nil
}

func (s *hpaHistoryManager) ingestEvent(ctx context.Context, clusterID uint, ev *corev1.Event) error {
	eventUID := string(ev.UID)
	if eventUID == "" {
		return nil
	}
	exists, err := s.historyRepo.ExistsByEventUID(ctx, eventUID)
	if err != nil || exists {
		return err
	}
	hpaName := ev.InvolvedObject.Name
	if ev.InvolvedObject.Kind != "HorizontalPodAutoscaler" || hpaName == "" {
		return nil
	}
	namespace := ev.Namespace
	if namespace == "" {
		namespace = ev.InvolvedObject.Namespace
	}
	client, err := s.cfgFactory.GetClientSet(ctx, clusterID)
	if err != nil {
		return err
	}
	hpa, err := client.AutoscalingV2().HorizontalPodAutoscalers(namespace).Get(ctx, hpaName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	oldReplicas, _ := s.historyRepo.GetLatestNewReplicas(ctx, clusterID, namespace, hpaName)
	newReplicas := hpa.Status.DesiredReplicas
	if newReplicas == 0 && hpa.Status.CurrentReplicas > 0 {
		newReplicas = hpa.Status.CurrentReplicas
	}
	msg := ev.Message
	if len(msg) > maxMessageLen {
		msg = msg[:maxMessageLen]
	}
	rec := &model.HPAScalingHistory{
		ClusterID:       int(clusterID),
		Namespace:       namespace,
		HPAName:         hpaName,
		ScaleTargetKind: hpa.Spec.ScaleTargetRef.Kind,
		ScaleTargetName: hpa.Spec.ScaleTargetRef.Name,
		OldReplicas:     oldReplicas,
		NewReplicas:     newReplicas,
		Reason:          ev.Reason,
		Message:         msg,
		EventUID:        eventUID,
		CreatedAt:       ev.LastTimestamp.Time,
	}
	if rec.CreatedAt.IsZero() {
		rec.CreatedAt = ev.EventTime.Time
	}
	if rec.CreatedAt.IsZero() {
		rec.CreatedAt = ev.CreationTimestamp.Time
	}
	return s.historyRepo.Create(ctx, rec)
}
