package model

import "time"

type Route struct {
	ID            int            `gorm:"type:bigint;primaryKey" json:"id"`
	WebhookID     int            `gorm:"type:bigint;not null;uniqueIndex" json:"webhookId"`
	Enabled       *bool          `gorm:"type:tinyint(1);default:1" json:"enabled"`
	Priority      int            `gorm:"type:int;not null;default:100" json:"priority"`
	RouteContinue *bool          `gorm:"column:route_continue;type:tinyint(1);default:0" json:"continue"`
	Matchers      []RouteMatcher `gorm:"foreignKey:RouteID" json:"matchers,omitempty"`
	CreatedAt     *time.Time     `gorm:"datetime(3)" json:"createdAt,omitempty"`
	UpdatedAt     *time.Time     `gorm:"datetime(3)" json:"updatedAt,omitempty"`
}

func (Route) TableName() string {
	return "prom_lens_alert_route"
}

type RouteMatcher struct {
	ID        int        `gorm:"type:bigint;primaryKey" json:"id"`
	RouteID   int        `gorm:"type:bigint;not null" json:"routeId"`
	Label     string     `gorm:"type:varchar(64);not null" json:"label"`
	Operator  string     `gorm:"type:varchar(8);not null;default:=" json:"operator"`
	Value     string     `gorm:"type:varchar(255);not null" json:"value"`
	SortOrder int        `gorm:"type:int;not null;default:0" json:"sortOrder"`
	CreatedAt *time.Time `gorm:"datetime(3)" json:"createdAt,omitempty"`
}

func (RouteMatcher) TableName() string {
	return "prom_lens_alert_route_matcher"
}
