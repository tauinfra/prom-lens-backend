package model

import (
	"time"

	authnmodel "valyria-backend/internal/apps/authn/model"
)

// K8s 同步状态
const (
	SyncStatusSynced  = "synced"       // 已同步至 K8s
	SyncStatusPending = "pending_sync" // 待同步
	SyncStatusFailed  = "failed"       // 同步失败，待重试
	SyncStatusDead    = "dead"         // 已达最大重试，不再自动重试
)

type Permission struct {
	ID          uint             `gorm:"primaryKey" json:"id,omitempty"`
	UserID      uint             `gorm:"uniqueIndex:idx_perm_user_cluster_ns_role;index;not null" json:"userID,omitempty"`
	ClusterID   uint             `gorm:"uniqueIndex:idx_perm_user_cluster_ns_role;index;not null" json:"clusterID,omitempty"`
	Namespace   string           `gorm:"uniqueIndex:idx_perm_user_cluster_ns_role;size:128;not null;default:'*'" json:"namespace,omitempty"`
	User        *authnmodel.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Cluster     *Cluster         `gorm:"foreignKey:ClusterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Role        string           `gorm:"uniqueIndex:idx_perm_user_cluster_ns_role;size:64;not null" json:"role,omitempty"`
	Creator     string           `gorm:"type:varchar(64);comment:'创建人'" json:"creator,omitempty"`
	CreatedAt   time.Time        `gorm:"type:datetime(3);comment:'创建时间'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time        `gorm:"type:datetime(3);comment:'更新时间'" json:"updatedAt,omitempty"`
	SyncStatus  string           `gorm:"type:varchar(32);default:'synced';comment:'K8s 同步状态'" json:"syncStatus,omitempty"`
	RetryCount  int              `gorm:"default:0;comment:'同步重试次数'" json:"retryCount,omitempty"`
	LastError   string           `gorm:"type:varchar(512);comment:'最近一次同步错误'" json:"lastError,omitempty"`
	LastErrorAt *time.Time       `gorm:"type:datetime(3);comment:'最近一次错误时间'" json:"lastErrorAt,omitempty"`
	LastSyncAt  *time.Time       `gorm:"type:datetime(3);comment:'最近一次同步成功时间'" json:"lastSyncAt,omitempty"`
	NextRetryAt *time.Time       `gorm:"type:datetime(3);comment:'下次重试时间'" json:"nextRetryAt,omitempty"`
}

func (Permission) TableName() string {
	return "valyria_kubernetes_permission"
}
