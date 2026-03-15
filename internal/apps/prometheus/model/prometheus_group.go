package model

import "time"

// Group 组
// 方案一: 使用 thanos ruler 维护告警规则
// 方案二: 使用 prometheus 维护告警规则，Group 需要关联 Server (不推荐，需要连接多个集群，存在安全隐患)
type Group struct {
	ID          int        `gorm:"type:bigint;primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(64);not null;unique;comment:'规则组名称'" json:"name,omitempty"`
	Type        string     `gorm:"type:varchar(16);not null;check:type IN ('ALERTING RULES', 'ALERTING RECORDS');comment:'组类型'" json:"type,omitempty"` // 规则组类型
	Description string     `gorm:"type:varchar(64);comment:'规则组描述'" json:"description,omitempty"`
	Rules       []Rule     `gorm:"<-:false" json:"rules,omitempty"`   // GORM 允许读，禁止写
	Records     []Record   `gorm:"<-:false" json:"records,omitempty"` // GORM 允许读，禁止写
	CreatedAt   *time.Time `gorm:"datetime(3);comment:'创建时间'" json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `gorm:"datetime(3);comment:'更新时间'" json:"updatedAt,omitempty"`
}

func (Group) TableName() string {
	return "valyria_prometheus_group"
}
