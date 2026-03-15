package middleware

import (
	"net/url"
	"strings"
)

// HasWhiteList 检查 URL 是否在白名单中（共享方法）
func HasWhiteList(url *url.URL, whiteList []string) bool {
	if len(whiteList) == 0 {
		return false
	}
	path := strings.ToLower(url.Path)
	for _, p := range whiteList {
		if strings.HasPrefix(path, strings.ToLower(p)) {
			return true
		}
	}
	return false
}
