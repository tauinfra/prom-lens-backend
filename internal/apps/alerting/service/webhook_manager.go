package service

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"strings"
	"time"

	"prom-lens-backend/internal/apps/alerting/dto"
	"prom-lens-backend/internal/apps/alerting/executor"
	"prom-lens-backend/internal/apps/alerting/model"
	"prom-lens-backend/internal/apps/alerting/repository"
	"prom-lens-backend/internal/apps/alerting/request"
	"prom-lens-backend/internal/core/logger"
	pg "prom-lens-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type WebhookManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.WebhookDTO, pg.Pagination, error)
	Get(ctx context.Context, id int) (dto.WebhookDTO, error)
	Create(ctx context.Context, req *request.CreateWebhookRequest) (dto.WebhookDTO, AMSyncResult, error)
	Update(ctx context.Context, id int, req *request.UpdateWebhookRequest) (AMSyncResult, error)
	Delete(ctx context.Context, id int) (AMSyncResult, error)
	Verify(ctx context.Context, req *request.VerifyWebhookRequest) error
	SyncReceivers(ctx context.Context) AMSyncResult
}

type webhookManager struct {
	repo   repository.WebhookRepository
	lark   *LarkClient
	amSync executor.AlertmanagerSyncer
}

func NewWebhookManager(repo repository.WebhookRepository, lark *LarkClient, amSync executor.AlertmanagerSyncer) WebhookManager {
	return &webhookManager{repo: repo, lark: lark, amSync: amSync}
}

func (s *webhookManager) List(ctx context.Context, params pg.QueryParams) ([]dto.WebhookDTO, pg.Pagination, error) {
	data, pagination, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, pg.Pagination{}, err
	}
	result := make([]dto.WebhookDTO, 0, len(data))
	for _, item := range data {
		result = append(result, dto.ToWebhookDTO(item))
	}
	return result, pagination, nil
}

func (s *webhookManager) Get(ctx context.Context, id int) (dto.WebhookDTO, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil {
		return dto.WebhookDTO{}, err
	}
	return dto.ToWebhookDTO(item), nil
}

func (s *webhookManager) Create(ctx context.Context, req *request.CreateWebhookRequest) (dto.WebhookDTO, AMSyncResult, error) {
	if err := ValidateLarkWebhookURL(req.URL); err != nil {
		return dto.WebhookDTO{}, AMSyncResult{}, err
	}
	route, matchers, err := buildRouteFromInput(req.Route)
	if err != nil {
		return dto.WebhookDTO{}, AMSyncResult{}, err
	}
	token, err := generateCallbackToken()
	if err != nil {
		return dto.WebhookDTO{}, AMSyncResult{}, err
	}
	enabled := boolDefault(req.Enabled, true)
	item := &model.Webhook{
		Name:          strings.TrimSpace(req.Name),
		URL:           strings.TrimSpace(req.URL),
		Description:   req.Description,
		Enabled:       &enabled,
		CallbackToken: token,
	}
	if err := s.repo.Create(ctx, item, route, matchers); err != nil {
		return dto.WebhookDTO{}, AMSyncResult{}, err
	}
	out, err := s.Get(ctx, item.ID)
	if err != nil {
		return dto.WebhookDTO{}, AMSyncResult{}, err
	}
	amSync := runAMSync(s.amSync, ctx)
	if !amSync.Success && executor.IsAlertmanagerSyncConfigured() {
		if delErr := s.repo.Delete(ctx, item.ID); delErr != nil {
			logger.Errorf("[am-sync] rollback webhook id=%d failed err=%v", item.ID, delErr)
		}
		return dto.WebhookDTO{}, amSync, fmt.Errorf("alertmanager 同步失败: %s", amSync.Msg)
	}
	return out, amSync, nil
}

func (s *webhookManager) Update(ctx context.Context, id int, req *request.UpdateWebhookRequest) (AMSyncResult, error) {
	data := &model.Webhook{}
	if req.Name != nil {
		data.Name = strings.TrimSpace(*req.Name)
	}
	if req.URL != nil {
		if err := ValidateLarkWebhookURL(*req.URL); err != nil {
			return AMSyncResult{}, err
		}
		data.URL = strings.TrimSpace(*req.URL)
	}
	if req.Description != nil {
		data.Description = *req.Description
	}
	if req.Enabled != nil {
		data.Enabled = req.Enabled
	}

	var route *model.Route
	var matchers []model.RouteMatcher
	replaceRoute := req.Route != nil
	if replaceRoute {
		var err error
		route, matchers, err = buildRouteFromInput(req.Route)
		if err != nil {
			return AMSyncResult{}, err
		}
	}
	if err := s.repo.Update(ctx, id, data, route, matchers, replaceRoute); err != nil {
		return AMSyncResult{}, err
	}
	return runAMSync(s.amSync, ctx), nil
}

func (s *webhookManager) Delete(ctx context.Context, id int) (AMSyncResult, error) {
	if err := s.repo.Delete(ctx, id); err != nil {
		return AMSyncResult{}, err
	}
	return runAMSync(s.amSync, ctx), nil
}

func (s *webhookManager) SyncReceivers(ctx context.Context) AMSyncResult {
	return runAMSync(s.amSync, ctx)
}

func (s *webhookManager) Verify(ctx context.Context, req *request.VerifyWebhookRequest) error {
	_ = ctx
	if err := ValidateLarkWebhookURL(req.URL); err != nil {
		return err
	}
	url := strings.TrimSpace(req.URL)
	testMsg := buildTestAlertMessage(req.Name)
	larkMsg, err := s.lark.GenerateCard(testMsg)
	if err != nil {
		return err
	}
	if err := s.lark.Send(url, larkMsg); err != nil {
		logger.Errorf("[alert-notify] verify failed url=%s err=%v", url, err)
		return err
	}
	logger.Infof("[alert-notify] verify success url=%s", url)
	return nil
}

func buildTestAlertMessage(name string) model.AlertManagerMessage {
	channelName := strings.TrimSpace(name)
	if channelName == "" {
		channelName = "prom-lens-test"
	}
	now := time.Now().UTC().Format(time.RFC3339)
	return model.AlertManagerMessage{
		Status:   "firing",
		Receiver: channelName,
		Alerts: []model.Alert{
			{
				Status: "firing",
				Labels: map[string]string{
					"alertname": "PromLensWebhookTest",
					"severity":  "info",
				},
				Annotations: map[string]string{
					"summary":     "Prom Lens Webhook 测试消息",
					"description": "这是一条来自 Prom Lens 的测试通知。若收到此消息，说明 Webhook 地址配置正常。",
				},
				StartsAt: now,
			},
		},
	}
}

func boolDefault(v *bool, def bool) bool {
	if v == nil {
		return def
	}
	return *v
}

type NotifyManager interface {
	ReceiveAlert(ctx context.Context, channel, token string, msg *model.AlertManagerMessage) error
}

type notifyManager struct {
	webhooks repository.WebhookRepository
	lark     *LarkClient
}

func NewNotifyManager(webhooks repository.WebhookRepository, lark *LarkClient) NotifyManager {
	return &notifyManager{webhooks: webhooks, lark: lark}
}

func (s *notifyManager) ReceiveAlert(ctx context.Context, channel, token string, msg *model.AlertManagerMessage) error {
	webhook, err := s.webhooks.GetByName(ctx, channel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("lark channel %q does not exist", channel)
		}
		return err
	}
	if subtle.ConstantTimeCompare([]byte(strings.TrimSpace(token)), []byte(webhook.CallbackToken)) != 1 {
		return fmt.Errorf("invalid webhook token")
	}
	if webhook.Enabled != nil && !*webhook.Enabled {
		return fmt.Errorf("lark channel %q is disabled", channel)
	}

	larkMsg, err := s.lark.GenerateCard(*msg)
	if err != nil {
		return err
	}
	if err := s.lark.Send(webhook.URL, larkMsg); err != nil {
		logger.Errorf("[alert-notify] failed channel=%s err=%v", channel, err)
		return err
	}
	logger.Infof("[alert-notify] success channel=%s status=%s alerts=%d", channel, msg.Status, len(msg.Alerts))
	return nil
}
