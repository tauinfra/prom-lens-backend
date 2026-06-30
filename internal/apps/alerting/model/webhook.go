package model

import "time"

type Webhook struct {
	ID            int        `gorm:"type:bigint;primaryKey" json:"id"`
	Name          string     `gorm:"type:varchar(64);not null;unique;comment:'通道名'" json:"name"`
	URL           string     `gorm:"type:varchar(512);not null;comment:'Lark webhook URL'" json:"url"`
	Description   string     `gorm:"type:varchar(255);comment:'描述'" json:"description,omitempty"`
	Enabled       *bool      `gorm:"type:tinyint(1);default:1;comment:'是否启用'" json:"enabled"`
	CallbackToken string     `gorm:"type:varchar(64);not null;comment:'Alertmanager 回调鉴权 token'" json:"callbackToken,omitempty"`
	Route         *Route     `gorm:"foreignKey:WebhookID" json:"route,omitempty"`
	CreatedAt     *time.Time `gorm:"datetime(3);comment:'创建时间'" json:"createdAt"`
	UpdatedAt     *time.Time `gorm:"datetime(3);comment:'更新时间'" json:"updatedAt"`
}

func (Webhook) TableName() string {
	return "prom_lens_alert_webhook"
}
