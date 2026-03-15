package model

import (
	"time"
)

// HPAScalingHistory HPA 扩缩容历史（由 Event 同步或 Watcher 写入）
type HPAScalingHistory struct {
	ID                int64     `gorm:"type:bigint;primaryKey;autoIncrement;comment:'主键'"     json:"id,omitempty"`
	ClusterID         int       `gorm:"type:int;not null;index:idx_hpa_history_lookup;comment:'集群 ID'" json:"clusterId,omitempty"`
	Namespace         string    `gorm:"type:varchar(128);not null;index:idx_hpa_history_lookup;comment:'命名空间'" json:"namespace,omitempty"`
	HPAName           string    `gorm:"type:varchar(256);not null;index:idx_hpa_history_lookup;comment:'HPA 名称'" json:"hpaName,omitempty"`
	ScaleTargetKind   string    `gorm:"type:varchar(64);comment:'缩放目标 Kind'"              json:"scaleTargetKind,omitempty"`
	ScaleTargetName   string    `gorm:"type:varchar(256);comment:'缩放目标名称'"              json:"scaleTargetName,omitempty"`
	OldReplicas       int32     `gorm:"type:int;comment:'变更前副本数'"                      json:"oldReplicas,omitempty"`
	NewReplicas       int32     `gorm:"type:int;not null;comment:'变更后副本数'"            json:"newReplicas,omitempty"`
	Reason            string    `gorm:"type:varchar(128);comment:'原因(Event Reason)'"      json:"reason,omitempty"`
	Message           string    `gorm:"type:varchar(512);comment:'消息(Event Message 截断)'" json:"message,omitempty"`
	EventUID          string    `gorm:"type:varchar(64);uniqueIndex:uk_event_uid;comment:'K8s Event UID 去重'" json:"eventUid,omitempty"`
	CreatedAt         time.Time `gorm:"type:datetime(3);index;comment:'发生时间'"            json:"createdAt,omitempty"`
}

func (HPAScalingHistory) TableName() string {
	return "valyria_hpa_scaling_history"
}
