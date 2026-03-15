package controller

import (
	"net/http"
	"valyria-backend/internal/apps/dragon/request"
	"valyria-backend/internal/apps/dragon/service"

	"github.com/gin-gonic/gin"
)

// ReleaseReportController 发布报表
type ReleaseReportController struct {
	report service.ReleaseReportManager
}

func NewReleaseReportController(report service.ReleaseReportManager) *ReleaseReportController {
	return &ReleaseReportController{report: report}
}

// Summary 发布统计报表汇总
// GET /dragon/reports/summary?from=2026-01&to=2026-02
func (c *ReleaseReportController) Summary(ctx *gin.Context) {
	var req request.ReleaseReportRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	data, err := c.report.Summary(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}
