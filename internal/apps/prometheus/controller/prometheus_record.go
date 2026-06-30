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

// RecordController 控制器结构体
type RecordController struct {
	record service.RecordManager // 使用服务接口
}

func NewRecordController(record service.RecordManager) *RecordController {
	return &RecordController{record: record}
}

func (c *RecordController) List(ctx *gin.Context) {
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
	data, pagination, err := c.record.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *RecordController) Get(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "recordID")
	if !ok {
		return
	}
	id := int(idU)
	data, err := c.record.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *RecordController) Create(ctx *gin.Context) {
	groupIDU, ok := ginhelper.RequireUintParam(ctx, "groupID")
	if !ok {
		return
	}
	groupID := int(groupIDU)
	var req request.CreateRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.record.Create(ctx, groupID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *RecordController) Update(ctx *gin.Context) {
	recordIDU, ok := ginhelper.RequireUintParam(ctx, "recordID")
	if !ok {
		return
	}
	var req request.UpdateRecordRequest
	recordID := int(recordIDU)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.record.Update(ctx, recordID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *RecordController) Delete(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "recordID")
	if !ok {
		return
	}
	id := int(idU)
	if err := c.record.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
