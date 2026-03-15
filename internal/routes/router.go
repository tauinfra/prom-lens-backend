package routes

import (
	"valyria-backend/internal/pkg/di"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, provider *di.Provider) *gin.Engine {
	r := gin.Default()
	// 全局中间件（注册顺序：先全局 -> 后局部）
	r.Use(
		provider.AuditManager.AuditLogin(db),     // 登录审计中间件
		provider.AuditManager.AuditOperation(db), // 操作审计中间件
		provider.JWTManager.Cors(),               // 跨域中间件
		provider.JWTManager.RequireAuth(),        // 认证中间件
	)

	api := r.Group("/api/v1")

	DashboardRouters(api, provider)      // 仪表盘（独立模块）
	AuthnRouters(api, provider)         // 认证系统
	AuditRouters(api, provider)          // 审计日志
	KubernetesRouters(api, db, provider) // 容器管理
	PrometheusRouters(api, provider)     // 监控平台
	DragonRouters(api, provider)         // 发布平台
	return r
}
