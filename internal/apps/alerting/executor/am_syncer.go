package executor

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"prom-lens-backend/internal/apps/alerting/model"
	"prom-lens-backend/internal/apps/alerting/repository"
	"prom-lens-backend/internal/core/logger"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type AlertmanagerSyncer interface {
	SyncAll(ctx context.Context) error
}

type alertmanagerSyncer struct {
	webhooks repository.WebhookRepository
}

func NewAlertmanagerSyncer(webhooks repository.WebhookRepository) AlertmanagerSyncer {
	return &alertmanagerSyncer{webhooks: webhooks}
}

func (s *alertmanagerSyncer) SyncAll(ctx context.Context) error {
	ns, cmName, key, baseURL, defaultReceiver, ok := getAlertmanagerConfig()
	if !ok {
		return fmt.Errorf("alertmanager sync is not configured (need base_url and alerting.alertmanager namespace/configmap/config_key)")
	}

	webhooks, err := s.webhooks.ListAllWithRoutes(ctx)
	if err != nil {
		logger.Errorf("[am-sync] list webhooks failed err=%v", err)
		return err
	}

	managed := make(map[string]struct{}, len(webhooks))
	for _, w := range webhooks {
		managed[w.Name] = struct{}{}
	}

	client, err := clientSet()
	if err != nil {
		logger.Errorf("[am-sync] k8s client failed err=%v", err)
		return err
	}

	cm, err := client.CoreV1().ConfigMaps(ns).Get(ctx, cmName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		merged, merr := mergeAlertmanagerYAML(nil, webhooks, managed, baseURL, defaultReceiver)
		if merr != nil {
			logger.Errorf("[am-sync] merge yaml for new configmap failed err=%v", merr)
			return merr
		}
		cm = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cmName,
				Namespace: ns,
			},
			Data: map[string]string{key: string(merged)},
		}
		if _, err = client.CoreV1().ConfigMaps(ns).Create(ctx, cm, metav1.CreateOptions{}); err != nil {
			logger.Errorf("[am-sync] create configmap ns=%s name=%s err=%v", ns, cmName, err)
			return fmt.Errorf("create alertmanager configmap %s/%s: %w", ns, cmName, err)
		}
		logger.Infof("[am-sync] created configmap ns=%s name=%s key=%s receivers=%d routes=%d",
			ns, cmName, key, countEnabledReceivers(webhooks), countEnabledRoutes(webhooks))
		return nil
	}
	if err != nil {
		logger.Errorf("[am-sync] get configmap ns=%s name=%s err=%v", ns, cmName, err)
		return fmt.Errorf("get alertmanager configmap %s/%s: %w", ns, cmName, err)
	}
	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	merged, err := mergeAlertmanagerYAML([]byte(cm.Data[key]), webhooks, managed, baseURL, defaultReceiver)
	if err != nil {
		logger.Errorf("[am-sync] merge yaml failed err=%v", err)
		return err
	}

	cm.Data[key] = string(merged)
	if _, err = client.CoreV1().ConfigMaps(ns).Update(ctx, cm, metav1.UpdateOptions{}); err != nil {
		logger.Errorf("[am-sync] update configmap ns=%s name=%s err=%v", ns, cmName, err)
		return err
	}

	logger.Infof("[am-sync] success ns=%s configmap=%s key=%s receivers=%d routes=%d",
		ns, cmName, key, countEnabledReceivers(webhooks), countEnabledRoutes(webhooks))
	return nil
}

func mergeAlertmanagerYAML(existing []byte, webhooks []model.Webhook, managed map[string]struct{}, baseURL, defaultReceiver string) ([]byte, error) {
	cfg := map[string]interface{}{}
	if len(strings.TrimSpace(string(existing))) > 0 {
		if err := yaml.Unmarshal(existing, &cfg); err != nil {
			return nil, fmt.Errorf("parse alertmanager config: %w", err)
		}
	}

	receivers := filterReceivers(cfg["receivers"], managed)
	receivers = append(receivers, buildReceivers(webhooks, baseURL)...)

	routeRoot := getRouteMap(cfg)
	manualRoutes := filterRoutes(routeRoot["routes"], managed)
	promRoutes := buildRoutes(webhooks)
	routeRoot["routes"] = append(promRoutes, manualRoutes...)
	if defaultReceiver != "" {
		routeRoot["receiver"] = defaultReceiver
	}
	cfg["route"] = routeRoot
	cfg["receivers"] = receivers

	return yaml.Marshal(cfg)
}

func buildReceivers(webhooks []model.Webhook, baseURL string) []map[string]interface{} {
	out := make([]map[string]interface{}, 0)
	for _, w := range webhooks {
		if w.Enabled != nil && !*w.Enabled {
			continue
		}
		url := fmt.Sprintf("%s/api/v1/alerting/webhook/%s", baseURL, w.Name)
		out = append(out, map[string]interface{}{
			"name": w.Name,
			"webhook_configs": []map[string]interface{}{
				{
					"url":            url,
					"send_resolved":  true,
					"http_config": map[string]interface{}{
						"bearer_token": w.CallbackToken,
					},
				},
			},
		})
	}
	return out
}

type routeEntry struct {
	priority int
	data     map[string]interface{}
}

func buildRoutes(webhooks []model.Webhook) []map[string]interface{} {
	entries := make([]routeEntry, 0)
	for _, w := range webhooks {
		if w.Enabled != nil && !*w.Enabled {
			continue
		}
		if w.Route == nil || w.Route.Enabled != nil && !*w.Route.Enabled {
			continue
		}
		if len(w.Route.Matchers) == 0 {
			continue
		}
		matchers := make([]string, 0, len(w.Route.Matchers))
		for _, m := range w.Route.Matchers {
			matchers = append(matchers, formatAlertmanagerMatcher(m.Label, m.Operator, m.Value))
		}
		cont := false
		if w.Route.RouteContinue != nil {
			cont = *w.Route.RouteContinue
		}
		entries = append(entries, routeEntry{
			priority: w.Route.Priority,
			data: map[string]interface{}{
				"receiver": w.Name,
				"matchers": matchers,
				"continue": cont,
			},
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].priority != entries[j].priority {
			return entries[i].priority < entries[j].priority
		}
		return false
	})
	out := make([]map[string]interface{}, 0, len(entries))
	for _, e := range entries {
		out = append(out, e.data)
	}
	return out
}

func filterReceivers(raw interface{}, managed map[string]struct{}) []map[string]interface{} {
	items, ok := raw.([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := m["name"].(string)
		if _, isManaged := managed[name]; isManaged {
			continue
		}
		out = append(out, m)
	}
	return out
}

func filterRoutes(raw interface{}, managed map[string]struct{}) []map[string]interface{} {
	items, ok := raw.([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		recv, _ := m["receiver"].(string)
		if _, isManaged := managed[recv]; isManaged {
			continue
		}
		out = append(out, m)
	}
	return out
}

func getRouteMap(cfg map[string]interface{}) map[string]interface{} {
	raw, ok := cfg["route"]
	if !ok {
		return map[string]interface{}{}
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	return m
}

func countEnabledReceivers(webhooks []model.Webhook) int {
	n := 0
	for _, w := range webhooks {
		if w.Enabled == nil || *w.Enabled {
			n++
		}
	}
	return n
}

func countEnabledRoutes(webhooks []model.Webhook) int {
	n := 0
	for _, w := range webhooks {
		if w.Enabled != nil && !*w.Enabled {
			continue
		}
		if w.Route != nil && (w.Route.Enabled == nil || *w.Route.Enabled) && len(w.Route.Matchers) > 0 {
			n++
		}
	}
	return n
}
