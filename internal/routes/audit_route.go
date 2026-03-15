package routes

import (
	"valyria-backend/internal/pkg/di"

	"github.com/gin-gonic/gin"
)

func AuditRouters(rg *gin.RouterGroup, provider *di.Provider) {
	audit := rg.Group("/audit")
	{
		audit.GET("/auth-logs", provider.Audit.AuthLog.Controller.List)
		audit.GET("/logs", provider.Audit.AuditLog.Controller.List)
	}
}
