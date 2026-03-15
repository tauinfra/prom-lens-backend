package middleware

import (
	"net/http"
	"valyria-backend/internal/apps/authn/service"
	"valyria-backend/internal/core/config"
	"valyria-backend/internal/core/logger"

	"github.com/gin-gonic/gin"
)

type PermissionManager struct {
	whiteList []string
	auth      service.AuthManager
}

// NewPermissionManager 创建权限中间件管理器
func NewPermissionManager(cfg *config.Config, auth service.AuthManager) *PermissionManager {
	return &PermissionManager{
		whiteList: cfg.Auth.Whitelist,
		auth:      auth,
	}
}

func (m *PermissionManager) RequirePermission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if HasWhiteList(c.Request.URL, m.whiteList) {
			c.Next()
			return
		}
		isSuperuserVal, ok := c.Get("isSuperuser")
		if ok {
			if isSuperuser, ok := isSuperuserVal.(bool); ok && isSuperuser {
				c.Next()
				return
			}
		}
		uidVal, ok := c.Get("uid")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"msg":     "请求未授权，缺少认证信息",
			})
			return
		}
		uid, ok := uidVal.(int)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"msg":     "请求未授权，认证信息无效",
			})
			return
		}
		ok, err := m.auth.HasUserPermission(c, uint(uid), code)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    http.StatusInternalServerError,
				"msg":     "权限校验异常，请联系运维处理",
			})
			return
		}
		if !ok {
			logger.Infof("[权限校验] UserID=%d | Method=%s | URL=%s 访问被拒绝，缺少权限: '%s'",
				uid, c.Request.Method, c.Request.URL.Path, code)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"msg":     "访问被拒绝，您无权限执行该操作",
			})
			return
		}
		c.Next()
	}
}
