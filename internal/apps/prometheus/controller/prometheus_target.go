package controller

import (
	"net/http"
	"strconv"
	"prom-lens-backend/internal/apps/prometheus/request"
	"prom-lens-backend/internal/apps/prometheus/service"
	pg "prom-lens-backend/internal/core/pagination"
	"prom-lens-backend/internal/pkg/ginhelper"

	"github.com/gin-gonic/gin"
)

// TargetController 控制器结构体
type TargetController struct {
	target service.TargetManager // 使用服务接口
}

func NewTargetController(target service.TargetManager) *TargetController {
	return &TargetController{target: target}
}

func (c *TargetController) List(ctx *gin.Context) {
	groupIDU, ok := ginhelper.RequireUintParam(ctx, "groupID")
	if !ok {
		return
	}
	groupID := int(groupIDU)
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Filters:   map[string]interface{}{"group_id": groupID},
	}
	data, pagination, err := c.target.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *TargetController) Get(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "targetID")
	if !ok {
		return
	}
	id := int(idU)
	data, err := c.target.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *TargetController) Create(ctx *gin.Context) {
	groupIDU, ok := ginhelper.RequireUintParam(ctx, "groupID")
	if !ok {
		return
	}
	groupID := int(groupIDU)
	var req request.CreateTargetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.target.Create(ctx, groupID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *TargetController) Update(ctx *gin.Context) {
	targetIDU, ok := ginhelper.RequireUintParam(ctx, "targetID")
	if !ok {
		return
	}
	var req request.UpdateTargetRequest
	targetID := int(targetIDU)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.target.Update(ctx, targetID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *TargetController) Delete(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "targetID")
	if !ok {
		return
	}
	id := int(idU)
	if err := c.target.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
