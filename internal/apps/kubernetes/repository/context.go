package repository

import "context"

// contextKey 用于 context.WithValue 的键类型，避免与第三方 key 冲突
type contextKey string

// ImpersonateUsernameKey 存在 context 时，Repository 使用 GetClientSetAsUser(id, username) 做 Impersonate
const ImpersonateUsernameKey contextKey = "kubernetes.impersonate.username"

// ImpersonateUsername 从 context 读取要冒充的用户名，不存在或非 string 返回空
func ImpersonateUsername(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v := ctx.Value(ImpersonateUsernameKey)
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}
