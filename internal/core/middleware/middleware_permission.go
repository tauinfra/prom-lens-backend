package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"valyria-backend/internal/apps/auth/model"
	"valyria-backend/internal/core/configs"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/pkg/casbin/foundry"
	"valyria-backend/internal/pkg/casbin/k8s"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PermissionManager struct {
	db        *gorm.DB
	resources []string
	actions   []string
	whiteList []string
}

// NewPermissionManager 创建 JWT 管理器
func NewPermissionManager(db *gorm.DB, cfg *configs.Config) *PermissionManager {
	return &PermissionManager{
		db:        db,
		resources: cfg.K8s.Resources,
		actions:   cfg.K8s.Actions,
		whiteList: cfg.Auth.Whitelist,
	}
}

func (m *PermissionManager) hasWhiteList(url *url.URL) bool {
	path := strings.ToLower(url.Path)
	for _, p := range m.whiteList {
		if strings.HasPrefix(path, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

func (m *PermissionManager) hasSuperAdmin(tx *gorm.DB, username string) bool {
	var user model.User
	if err := tx.Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}
	return *user.IsSuperuser
}

func (m *PermissionManager) parseAction(method, url string) (action string) {
	parts := strings.Split(strings.Trim(url, "/"), "/")
	resource := ""
	resourceIndex := -1

	// 优先处理特殊 Action: logs/exec/rollout/restart 等
	for _, part := range parts {
		for _, action = range m.actions {
			if action == part {
				return action
			}
		}
	}

	// 倒序扫描找到最后一个匹配的资源
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if part == "" {
			continue
		}
		for _, resource = range m.resources {
			if part == resource {
				resourceIndex = i
				break
			}
		}
		if resource != "" {
			break
		}
	}

	// 根据 resource 后是否有 ID 来判断 list/get
	hasID := false
	if resourceIndex != -1 && resourceIndex+1 < len(parts) {
		nextPart := parts[resourceIndex+1]
		// 下一个部分不是资源名，认为是 ID
		isNextResource := false
		for _, r := range m.resources {
			if nextPart == r {
				isNextResource = true
				break
			}
		}
		hasID = nextPart != "" && !isNextResource
	}

	switch method {
	case "GET":
		if resource == "" {
			return "unknown"
		}
		if hasID {
			return "get"
		}
		return "list"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "unknown"
	}
}

func (m *PermissionManager) parseResourceURL(url string) (cluster, namespace, resource string) {
	parts := strings.Split(strings.Trim(url, "/"), "/")
	cluster = "*"
	namespace = "*"
	resource = "*"
	for i := 0; i < len(parts); i++ {
		// clusterID
		if parts[i] == "clusters" && i+1 < len(parts) {
			cluster = parts[i+1]
		}
		// namespace
		if parts[i] == "namespaces" && i+1 < len(parts) {
			namespace = parts[i+1]
		}
	}
	// resource 需要从后向前匹配
	for i := len(parts) - 1; i >= 0; i-- {
		for _, r := range m.resources {
			if parts[i] == r {
				resource = r
				return
			}
		}
	}
	return
}

func (m *PermissionManager) RBACPermission(tx *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if k8s.Enforcer == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"msg":     "Casbin not initialized",
			})
			return
		}

		if m.hasWhiteList(c.Request.URL) {
			c.Next() // 如果 url 在白名单，直接执行后续函数
			return
		}

		// 从上下文或 token 提取用户名
		username, exists := c.Get("username")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"msg":     "用户未登录",
			})
			return
		}

		// 超级管理
		if m.hasSuperAdmin(tx, username.(string)) {
			c.Next()
			return
		}

		cluster, namespace, resource := m.parseResourceURL(c.Request.URL.Path)
		act := m.parseAction(c.Request.Method, c.Request.URL.Path)
		// 权限判断
		ok, err := k8s.Enforcer.Enforce(username, cluster, namespace, resource, act)
		logger.Infof("current user request permission k8s: [%s|%s|%s|%s|%s] result: %v", username, cluster, namespace, resource, act, ok)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    http.StatusInternalServerError,
				"msg":     "权限校验错误，请联系运维人员处理！",
			})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"msg":     "访问被拒绝，您无权限执行该操作，请联系运维人员处理！",
			})
			return
		}
		c.Next()
	}
}

func (m *PermissionManager) ModuleAccess(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		var path = c.Request.URL.Path
		// 白名单｜无需授权
		if m.hasWhiteList(c.Request.URL) {
			c.Next() // 如果 url 在白名单，直接执行后续函数
			return
		}
		// 从上下文或 token 提取用户名
		username, exists := c.Get("username")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"msg":     "用户未登录",
			})
			return
		}

		// 超级管理
		if m.hasSuperAdmin(db, username.(string)) {
			c.Next()
			return
		}

		if err := db.Where("username = ?", username).Preload("Modules").First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    http.StatusInternalServerError,
				"msg":     "查询用户失败",
			})
			return
		}

		allowed := false
		for _, module := range user.Modules {
			if strings.HasPrefix(path, module.RoutePrefix) {
				allowed = true
				break
			}
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"msg":     "您无权限访问该模块，请联系运维人员处理！",
				"path":    path,
			})
			return
		}
		c.Next()
	}
}

// FoundryPermission 发布系统权限中间件
func (m *PermissionManager) FoundryPermission(tx *gorm.DB, resource, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if foundry.Enforcer == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"msg":     "Foundry module casbin not initialized.",
			})
			return
		}

		// 从上下文获取登录用户
		username, _ := c.Get("username")
		// 超级管理
		if m.hasSuperAdmin(tx, username.(string)) {
			if act == "list" {
				c.Set("accessIDs", []string{"*"})
			}
			c.Next()
			return
		}

		sub := username.(string)
		obj := resource

		projectID := c.Param("projectID")
		envID := c.Param("envID")
		pipelineID := c.Param("pipelineID")
		if act != "list" && act != "create" {
			if projectID != "" {
				obj = fmt.Sprintf("%s:%s", resource, projectID)
			}
			if projectID != "" && envID != "" {
				obj = fmt.Sprintf("%s:%s", resource, envID)
			}
			if projectID != "" && envID != "" && pipelineID != "" {
				obj = fmt.Sprintf("%s:%s", resource, pipelineID)
			}
		}
		ok, err := foundry.Enforcer.Enforce(sub, obj, act)
		logger.Infof("foundry permission check: [%s | %s | %s] => %v, err: %v", sub, obj, act, ok, err)
		// 报错或全部未命中 → 拒绝
		if err != nil || !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"msg":     "访问被拒绝，您无权限执行该操作，请联系运维人员处理！",
			})
			return
		}
		//
		if act == "list" {
			ids := []string{}
			policies, err := foundry.Enforcer.GetFilteredPolicy(0, sub)
			if err != nil || !ok {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"code":    http.StatusInternalServerError,
					"msg":     "获取权限策略异常，错误原因: " + err.Error(),
				})
				return
			}
			for _, policy := range policies {
				if policy[2] != "get" { // 只关注 get 权限
					continue
				}
				if policy[1] == resource+":*" {
					ids = []string{"*"} // 全量权限
					break
				}
				prefix := resource + ":"
				if len(policy[1]) > len(prefix) && policy[1][:len(prefix)] == prefix {
					ids = append(ids, policy[1][len(prefix):])
				}
			}
			c.Set("accessIDs", ids)
		}
		c.Next()
	}
}
