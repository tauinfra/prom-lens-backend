package controller

import (
	"net/http"
	"strconv"
	"strings"

	"prom-lens-backend/internal/apps/alerting/model"
	"prom-lens-backend/internal/apps/alerting/request"
	"prom-lens-backend/internal/apps/alerting/service"
	pg "prom-lens-backend/internal/core/pagination"
	"prom-lens-backend/internal/pkg/ginhelper"

	"github.com/gin-gonic/gin"
)

type WebhookController struct {
	webhook service.WebhookManager
}

func NewWebhookController(webhook service.WebhookManager) *WebhookController {
	return &WebhookController{webhook: webhook}
}

func (c *WebhookController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"),
	}
	data, pagination, err := c.webhook.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *WebhookController) Get(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "webhookID")
	if !ok {
		return
	}
	data, err := c.webhook.Get(ctx, int(idU))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *WebhookController) Create(ctx *gin.Context) {
	var req request.CreateWebhookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	data, amSync, err := c.webhook.Create(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error(), "amSync": amSync})
		return
	}
	if !amSync.Success {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": amSync.Msg, "data": data, "amSync": amSync})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "amSync": amSync})
}

func (c *WebhookController) Update(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "webhookID")
	if !ok {
		return
	}
	var req request.UpdateWebhookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	amSync, err := c.webhook.Update(ctx, int(idU), &req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error(), "amSync": amSync})
		return
	}
	if !amSync.Success {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": amSync.Msg, "amSync": amSync})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "amSync": amSync})
}

func (c *WebhookController) Delete(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "webhookID")
	if !ok {
		return
	}
	amSync, err := c.webhook.Delete(ctx, int(idU))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error(), "amSync": amSync})
		return
	}
	if !amSync.Success {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": amSync.Msg, "amSync": amSync})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "amSync": amSync})
}

func (c *WebhookController) SyncReceivers(ctx *gin.Context) {
	amSync := c.webhook.SyncReceivers(ctx)
	if !amSync.Success {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": amSync.Msg, "amSync": amSync})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "amSync": amSync})
}

func (c *WebhookController) Verify(ctx *gin.Context) {
	var req request.VerifyWebhookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.webhook.Verify(ctx, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

type NotifyController struct {
	notify service.NotifyManager
}

func NewNotifyController(notify service.NotifyManager) *NotifyController {
	return &NotifyController{notify: notify}
}

func extractBearerToken(ctx *gin.Context) string {
	authHeader := strings.TrimSpace(ctx.GetHeader("Authorization"))
	if len(authHeader) >= 7 && strings.EqualFold(authHeader[:7], "Bearer ") {
		return strings.TrimSpace(authHeader[7:])
	}
	return ""
}

func (c *NotifyController) Webhook(ctx *gin.Context) {
	channel := strings.TrimSpace(ctx.Param("channel"))
	if channel == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": -1, "error": "channel is required"})
		return
	}
	token := extractBearerToken(ctx)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": -1, "error": "Authorization Bearer token is required"})
		return
	}
	var alertMsg model.AlertManagerMessage
	if err := ctx.ShouldBindJSON(&alertMsg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": -1, "error": err.Error()})
		return
	}
	if err := c.notify.ReceiveAlert(ctx, channel, token, &alertMsg); err != nil {
		if strings.Contains(err.Error(), "invalid webhook token") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": -1, "error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "disabled") {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": -1, "error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": -1, "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "status": "success"})
}
