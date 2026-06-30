package dto

import (
	"time"

	"prom-lens-backend/internal/apps/alerting/model"
)

type WebhookDTO struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	URL           string     `json:"url"`
	Description   string     `json:"description"`
	Enabled       *bool      `json:"enabled"`
	CallbackToken string     `json:"callbackToken"`
	Route         *RouteDTO  `json:"route,omitempty"`
	CreatedAt     *time.Time `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
}

func ToWebhookDTO(item model.Webhook) WebhookDTO {
	return WebhookDTO{
		ID:            item.ID,
		Name:          item.Name,
		URL:           item.URL,
		Description:   item.Description,
		Enabled:       item.Enabled,
		CallbackToken: item.CallbackToken,
		Route:         ToRouteDTO(item.Route),
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}
}
