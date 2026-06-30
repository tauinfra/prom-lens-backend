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

// GroupController 控制器结构体
type GroupController struct {
	group service.GroupManager // 使用服务接口
}

func NewGroupController(group service.GroupManager) *GroupController {
	return &GroupController{group: group}
}

func (c *GroupController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"), // 搜索关键字
	}
	if groupType := ctx.Query("type"); groupType != "" {
		params.Filters = map[string]interface{}{"type": groupType}
	}
	data, pagination, err := c.group.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *GroupController) Get(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "groupID")
	if !ok {
		return
	}
	id := int(idU)
	data, err := c.group.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *GroupController) Create(ctx *gin.Context) {
	var req request.CreateGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.group.Create(ctx, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *GroupController) Update(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "groupID")
	if !ok {
		return
	}
	var req request.UpdateGroupRequest
	id := int(idU)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.group.Update(ctx, id, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *GroupController) Delete(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "groupID")
	if !ok {
		return
	}
	id := int(idU)
	if err := c.group.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
