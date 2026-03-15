package repository

import (
	"context"
	"strings"
	"valyria-backend/internal/apps/kubernetes/model"

	"gorm.io/gorm"
)

// HpaHistoryRepository HPA 扩缩容历史 DB 仓储（按 Event 同步落库）
type HpaHistoryRepository interface {
	Create(ctx context.Context, record *model.HPAScalingHistory) error
	ExistsByEventUID(ctx context.Context, eventUID string) (bool, error)
	GetLatestNewReplicas(ctx context.Context, clusterID uint, namespace, hpaName string) (int32, error)
	List(ctx context.Context, clusterID uint, namespace, hpaName string, page, size int) ([]model.HPAScalingHistory, int64, error)
}

type hpaHistoryRepository struct {
	db *gorm.DB
}

func NewHpaHistoryRepository(db *gorm.DB) HpaHistoryRepository {
	return &hpaHistoryRepository{db: db}
}

// Create 插入一条历史；若 event_uid 已存在（唯一约束）视为幂等成功，返回 nil
func (r *hpaHistoryRepository) Create(ctx context.Context, record *model.HPAScalingHistory) error {
	err := r.db.WithContext(ctx).Create(record).Error
	if err != nil && isDuplicateKeyError(err) {
		return nil
	}
	return err
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	// MySQL 1062, PostgreSQL 23505, SQLite UNIQUE constraint
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "duplicate") || strings.Contains(s, "1062") || strings.Contains(s, "23505") || strings.Contains(s, "unique constraint")
}

func (r *hpaHistoryRepository) ExistsByEventUID(ctx context.Context, eventUID string) (bool, error) {
	if eventUID == "" {
		return false, nil
	}
	var n int64
	err := r.db.WithContext(ctx).Model(&model.HPAScalingHistory{}).Where("event_uid = ?", eventUID).Limit(1).Count(&n).Error
	return n > 0, err
}

func (r *hpaHistoryRepository) GetLatestNewReplicas(ctx context.Context, clusterID uint, namespace, hpaName string) (int32, error) {
	var rec model.HPAScalingHistory
	err := r.db.WithContext(ctx).Model(&model.HPAScalingHistory{}).
		Where("cluster_id = ? AND namespace = ? AND hpa_name = ?", clusterID, namespace, hpaName).
		Order("created_at DESC").Limit(1).First(&rec).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return rec.NewReplicas, nil
}

func (r *hpaHistoryRepository) List(ctx context.Context, clusterID uint, namespace, hpaName string, page, size int) ([]model.HPAScalingHistory, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.HPAScalingHistory{}).
		Where("cluster_id = ? AND namespace = ? AND hpa_name = ?", clusterID, namespace, hpaName)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	var list []model.HPAScalingHistory
	err := q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}
