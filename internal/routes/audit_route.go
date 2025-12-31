package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	controller2 "valyria-backend/internal/apps/audit/controller"
)

func AuditRouters(tx *gorm.DB, r *gin.Engine) {
	// Audit Routes
	var authLog = controller2.SetupAuthLogController(tx)
	var auditLog = controller2.SetupAuditLogController(tx)
	AuditRouters := r.Group("/api/v1/audit")
	{
		AuditRouters.GET("/auth-logs", authLog.List)
	}
	{
		AuditRouters.GET("/logs", auditLog.List)
	}
}
