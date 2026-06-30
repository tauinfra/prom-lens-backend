package routes

import (
	"prom-lens-backend/internal/pkg/di"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, provider *di.Provider) *gin.Engine {
	r := gin.Default()
	registerHealthRoutes(r, db)
	// 全局中间件（注册顺序：先全局 -> 后局部）
	r.Use(
		provider.AuditManager.AuditLogin(db),     // 登录审计中间件
		provider.AuditManager.AuditOperation(db), // 操作审计中间件
		provider.JWTManager.Cors(),               // 跨域中间件
		provider.JWTManager.RequireAuth(),        // 认证中间件
	)

	api := r.Group("/api/v1")
	api.Use(provider.JWTManager.RequireSuperuserForDelete())

	AuthnRouters(api, provider)       // 认证系统
	AuditRouters(api, provider)       // 审计日志
	PrometheusRouters(api, provider)  // Prometheus 规则配置
	AlertingRouters(api, provider)    // 告警通知
	return r
}
