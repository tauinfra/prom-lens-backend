package service

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

var allowedLarkHosts = map[string]struct{}{
	"open.feishu.cn":     {},
	"open.larksuite.com": {},
}

const larkWebhookPathPrefix = "/open-apis/bot/v2/hook/"

// ValidateLarkWebhookURL 校验 Lark Webhook URL，防止 SSRF。
func ValidateLarkWebhookURL(raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fmt.Errorf("url is required")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid lark webhook url")
	}
	if u.Scheme != "https" {
		return fmt.Errorf("lark webhook url must use https")
	}
	if u.User != nil {
		return fmt.Errorf("invalid lark webhook url")
	}

	host := strings.ToLower(u.Hostname())
	if _, ok := allowedLarkHosts[host]; !ok {
		return fmt.Errorf("lark webhook url must use feishu or larksuite domain")
	}
	if ip := net.ParseIP(host); ip != nil {
		return fmt.Errorf("invalid lark webhook url")
	}
	if !strings.HasPrefix(u.Path, larkWebhookPathPrefix) || u.Path == larkWebhookPathPrefix {
		return fmt.Errorf("invalid lark webhook path")
	}
	if u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("invalid lark webhook url")
	}
	return nil
}
