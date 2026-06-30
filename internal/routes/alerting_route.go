package routes

import (
	"prom-lens-backend/internal/pkg/di"

	"github.com/gin-gonic/gin"
)

func AlertingRouters(rg *gin.RouterGroup, provider *di.Provider) {
	alerting := rg.Group("/alerting")
	{
		alerting.GET("/webhooks", provider.Alerting.Webhook.Controller.List)
		alerting.GET("/webhooks/:webhookID", provider.Alerting.Webhook.Controller.Get)
		alerting.POST("/webhooks", provider.Alerting.Webhook.Controller.Create)
		alerting.POST("/webhooks/verify", provider.Alerting.Webhook.Controller.Verify)
		alerting.PATCH("/webhooks/:webhookID", provider.Alerting.Webhook.Controller.Update)
		alerting.DELETE("/webhooks/:webhookID", provider.Alerting.Webhook.Controller.Delete)
		alerting.POST("/sync/receivers", provider.Alerting.Webhook.Controller.SyncReceivers)

		alerting.POST("/webhook/:channel", provider.Alerting.Notify.Controller.Webhook)
	}
}
