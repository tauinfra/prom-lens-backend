package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"prom-lens-backend/internal/apps/alerting/model"
)

type LarkClient struct {
	timeout time.Duration
	client  *http.Client
}

func NewLarkClient(timeout time.Duration) *LarkClient {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &LarkClient{
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *LarkClient) Send(webhookURL string, msg model.LarkMsg) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal lark message: %w", err)
	}
	resp, err := s.client.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("send lark message: %w", err)
	}
	if resp == nil {
		return fmt.Errorf("send lark message: empty response")
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("lark webhook status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func (s *LarkClient) GenerateCard(amMsg model.AlertManagerMessage) (model.LarkMsg, error) {
	if len(amMsg.Alerts) == 0 {
		return model.LarkMsg{}, fmt.Errorf("no alerts to send")
	}
	status := strings.ToLower(amMsg.Status)
	title := fmt.Sprintf("[%s: %d] Prometheus Alert", strings.ToUpper(amMsg.Status), len(amMsg.Alerts))
	firstSummary := strings.TrimSpace(amMsg.Alerts[0].Annotations["summary"])
	if firstSummary != "" {
		title = fmt.Sprintf("[%s: %d] %s", strings.ToUpper(amMsg.Status), len(amMsg.Alerts), firstSummary)
	}
	template := "blue"
	switch status {
	case "firing":
		template = "red"
	case "resolved":
		template = "green"
	}

	elements := make([]model.LarkCardElement, 0, len(amMsg.Alerts)*2)
	for i, alert := range amMsg.Alerts {
		if i > 0 {
			elements = append(elements, model.LarkCardElement{Tag: "hr"})
		}
		elements = append(elements, model.LarkCardElement{
			Tag: "div",
			Text: &model.LarkCardText{
				Tag:     "lark_md",
				Content: buildAlertMarkdown(alert, amMsg.Status),
			},
		})
	}

	return model.LarkMsg{
		MsgType: "interactive",
		Card: &model.LarkCard{
			Header: model.LarkCardHeader{
				Title: model.LarkCardText{
					Tag:     "plain_text",
					Content: title,
				},
				Template: template,
			},
			Elements: elements,
		},
	}, nil
}

func formatAlertTime(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t.Format("2006-01-02 15:04:05")
	}
	if t, err := time.Parse(time.RFC3339Nano, input); err == nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return input
}

func buildAlertMarkdown(alert model.Alert, overallStatus string) string {
	startsAt := formatAlertTime(alert.StartsAt)
	if startsAt == "" {
		startsAt = "N/A"
	}
	endsAt := formatAlertTime(alert.EndsAt)
	if endsAt == "" {
		endsAt = "N/A"
	}

	lines := []string{}
	if desc := strings.TrimSpace(alert.Annotations["description"]); desc != "" {
		lines = append(lines, fmt.Sprintf("**description:** %s", desc))
	}
	if len(alert.Labels) > 0 {
		keys := make([]string, 0, len(alert.Labels))
		for k := range alert.Labels {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			lines = append(lines, fmt.Sprintf("**%s: ** %s", k, alert.Labels[k]))
		}
	}
	lines = append(lines, fmt.Sprintf("**startsAt:** %s", startsAt))
	status := strings.ToLower(alert.Status)
	if status == "resolved" || strings.ToLower(overallStatus) == "resolved" {
		lines = append(lines, fmt.Sprintf("**endsAt:** %s", endsAt))
	}
	return strings.Join(lines, "\n")
}
