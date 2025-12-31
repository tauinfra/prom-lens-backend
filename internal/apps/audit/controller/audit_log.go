package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/audit/repository"
	"valyria-backend/internal/apps/audit/service"
	pg "valyria-backend/internal/core/pagination"
)

// AuditLogController 处理集群相关请求的控制器结构体
type AuditLogController struct {
	auditLog service.AuditLogService // 使用服务接口
}

// SetupAuditLogController 创建新的 AuditLogController 实例
func SetupAuditLogController(tx *gorm.DB) *AuditLogController {
	auditLog := service.NewAuditLogService(tx, repository.NewAuditLogRepository(tx))
	return &AuditLogController{auditLog: auditLog}
}

func (c *AuditLogController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))

	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
	}
	// 默认排序
	if params.SortBy == "" && params.SortOrder == "" {
		params.SortBy = "created_at"
		params.SortOrder = "desc"
	}
	// 查询
	data, pagination, err := c.auditLog.List(params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}
