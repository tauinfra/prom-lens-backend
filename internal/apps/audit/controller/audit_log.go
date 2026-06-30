package controller

import (
	"net/http"
	"strconv"
	"prom-lens-backend/internal/apps/audit/service"
	pg "prom-lens-backend/internal/core/pagination"

	"github.com/gin-gonic/gin"
)

// AuditLogController 处理审计日志请求的控制器结构体
type AuditLogController struct {
	auditLog service.AuditLogService // 使用服务接口
}

// NewAuditLogController 创建新的 AuditLogController 实例
func NewAuditLogController(auditLog service.AuditLogService) *AuditLogController {
	return &AuditLogController{auditLog: auditLog}
}

func (c *AuditLogController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"),
	}
	// 查询
	data, pagination, err := c.auditLog.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}
