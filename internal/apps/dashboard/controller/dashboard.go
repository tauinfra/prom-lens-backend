package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/dashboard/service"

	"github.com/gin-gonic/gin"
)

// DashboardController 仪表盘控制器
type DashboardController struct {
	dashboard service.DashboardManager
}

// NewDashboardController 创建仪表盘控制器
func NewDashboardController(dashboard service.DashboardManager) *DashboardController {
	return &DashboardController{dashboard: dashboard}
}

// Get 返回仪表盘聚合数据：集群/节点/命名空间/Pod；resources 预留；pipelines 为今日发布统计
func (c *DashboardController) Get(ctx *gin.Context) {
	data, err := c.dashboard.GetDashboard(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

// GetPipelinesTrend 流水线趋势，query: range=today|7d|30d
func (c *DashboardController) GetPipelinesTrend(ctx *gin.Context) {
	rangeParam := ctx.DefaultQuery("range", "today")
	if rangeParam != "today" && rangeParam != "7d" && rangeParam != "30d" {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": "range must be today, 7d or 30d"})
		return
	}
	data, err := c.dashboard.GetPipelinesTrend(ctx, rangeParam)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

// GetPipelinesProjects 按项目统计发布次数、pipeline 数、成功/失败/回滚次数
func (c *DashboardController) GetPipelinesProjects(ctx *gin.Context) {
	data, err := c.dashboard.GetPipelinesProjects(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

// GetPipelinesRecent 最近发布记录，query: env=prod&limit=10
func (c *DashboardController) GetPipelinesRecent(ctx *gin.Context) {
	env := ctx.Query("env")
	if env == "" {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": "env is required"})
		return
	}
	limit := 10
	if l := ctx.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}
	data, err := c.dashboard.GetPipelinesRecent(ctx, env, limit)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}
