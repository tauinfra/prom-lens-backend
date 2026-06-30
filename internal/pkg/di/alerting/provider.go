package alerting

import (
	"time"
	"prom-lens-backend/internal/apps/alerting/controller"
	"prom-lens-backend/internal/apps/alerting/executor"
	"prom-lens-backend/internal/apps/alerting/repository"
	"prom-lens-backend/internal/apps/alerting/service"
	"prom-lens-backend/internal/core/config"

	"gorm.io/gorm"
)

type Provider struct {
	Webhook *WebhookProvider
	Notify  *NotifyProvider
}

type WebhookProvider struct {
	Service    service.WebhookManager
	Controller *controller.WebhookController
}

type NotifyProvider struct {
	Service    service.NotifyManager
	Controller *controller.NotifyController
}

func NewAlertingProvider(db *gorm.DB, cfg *config.Config) *Provider {
	repo := repository.NewWebhookRepository(db)

	timeout := 10 * time.Second
	if cfg.Alerting.LarkTimeout > 0 {
		timeout = cfg.Alerting.LarkTimeout
	}
	larkClient := service.NewLarkClient(timeout)
	amSyncer := executor.NewAlertmanagerSyncer(repo)
	webhookSvc := service.NewWebhookManager(repo, larkClient, amSyncer)
	notifySvc := service.NewNotifyManager(repo, larkClient)

	return &Provider{
		Webhook: &WebhookProvider{
			Service:    webhookSvc,
			Controller: controller.NewWebhookController(webhookSvc),
		},
		Notify: &NotifyProvider{
			Service:    notifySvc,
			Controller: controller.NewNotifyController(notifySvc),
		},
	}
}
