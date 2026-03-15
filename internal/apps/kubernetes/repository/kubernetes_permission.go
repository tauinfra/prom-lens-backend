package repository

import (
	"context"
	"time"
	"valyria-backend/internal/apps/kubernetes/model"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	List(ctx context.Context, clusterID uint) ([]model.Permission, error)
	ListByUserID(ctx context.Context, userID uint) ([]model.Permission, error)
	ListPendingSync(ctx context.Context, limit int, maxRetry int) ([]model.Permission, error)
	Create(ctx context.Context, data *model.Permission) error
	Delete(ctx context.Context, id uint) error
	Get(ctx context.Context, id uint) (model.Permission, error)
	GetUsernameByUserID(ctx context.Context, userID uint) (string, error)
	UpdateSyncSuccess(ctx context.Context, id uint, lastSyncAt time.Time) error
	UpdateSyncFailure(ctx context.Context, id uint, retryCount int, lastError string, nextRetryAt *time.Time, maxRetry int) error
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) List(ctx context.Context, clusterID uint) ([]model.Permission, error) {
	var list []model.Permission
	q := r.db.WithContext(ctx).Preload("User").Preload("Cluster")
	if clusterID > 0 {
		q = q.Where("cluster_id = ?", clusterID)
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *permissionRepository) ListByUserID(ctx context.Context, userID uint) ([]model.Permission, error) {
	var list []model.Permission
	err := r.db.WithContext(ctx).Preload("User").Preload("Cluster").Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

// defaultMaxRetry 与 config permission_worker 默认一致，仅当调用方传入 <=0 时使用
const defaultMaxRetry = 5

// ListPendingSync 查询待同步或失败且到达下次重试时间的记录（幂等：仅 next_retry_at <= now 或为 null）
// maxRetry 若 <=0 使用 defaultMaxRetry
func (r *permissionRepository) ListPendingSync(ctx context.Context, limit int, maxRetry int) ([]model.Permission, error) {
	if limit <= 0 {
		limit = 100
	}
	threshold := maxRetry
	if threshold <= 0 {
		threshold = defaultMaxRetry
	}
	var list []model.Permission
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("sync_status IN ? AND retry_count < ?", []string{model.SyncStatusPending, model.SyncStatusFailed}, threshold).
		Where("next_retry_at IS NULL OR next_retry_at <= ?", now).
		Order("next_retry_at ASC").
		Limit(limit).
		Find(&list).Error
	return list, err
}

func (r *permissionRepository) Create(ctx context.Context, data *model.Permission) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *permissionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Permission{}, id).Error
}

func (r *permissionRepository) Get(ctx context.Context, id uint) (model.Permission, error) {
	var m model.Permission
	err := r.db.WithContext(ctx).Preload("User").Preload("Cluster").First(&m, id).Error
	return m, err
}

// GetUsernameByUserID 从 authn 用户表查用户名，用于 K8s RBAC 同步
func (r *permissionRepository) GetUsernameByUserID(ctx context.Context, userID uint) (string, error) {
	var name string
	err := r.db.WithContext(ctx).Table("valyria_authn_user").Where("id = ?", userID).Limit(1).Select("username").Scan(&name).Error
	return name, err
}

func (r *permissionRepository) UpdateSyncSuccess(ctx context.Context, id uint, lastSyncAt time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Permission{}).Where("id = ?", id).Updates(map[string]interface{}{
		"sync_status":   model.SyncStatusSynced,
		"last_sync_at":  lastSyncAt,
		"retry_count":   0,
		"last_error":    "",
		"last_error_at": nil,
		"next_retry_at": nil,
	}).Error
}

// UpdateSyncFailure 更新失败状态；若 retryCount >= maxRetry 则置为 dead 且 next_retry_at 清空；maxRetry<=0 用 defaultMaxRetry
func (r *permissionRepository) UpdateSyncFailure(ctx context.Context, id uint, retryCount int, lastError string, nextRetryAt *time.Time, maxRetry int) error {
	threshold := maxRetry
	if threshold <= 0 {
		threshold = defaultMaxRetry
	}
	now := time.Now()
	upd := map[string]interface{}{
		"retry_count":   retryCount,
		"last_error":    lastError,
		"last_error_at": now,
	}
	if retryCount >= threshold {
		upd["sync_status"] = model.SyncStatusDead
		upd["next_retry_at"] = nil
	} else {
		upd["sync_status"] = model.SyncStatusFailed
		upd["next_retry_at"] = nextRetryAt
	}
	return r.db.WithContext(ctx).Model(&model.Permission{}).Where("id = ?", id).Updates(upd).Error
}
